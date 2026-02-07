package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/yaml"
)

var (
	version   = "dev"
	sigHangup = make(chan os.Signal, 1)
)

func main() {
	cmd := &cli.Command{
		Name:    "showbridge",
		Usage:   "Simple protocol router /s",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./config.yaml",
				Usage: "path to config file",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Usage: "set log level",
				Validator: func(level string) error {
					levels := []string{"debug", "info", "warn", "error"}
					if !slices.Contains(levels, level) {
						return fmt.Errorf("unknown log level: %s", level)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "log-format",
				Value: "text",
				Usage: "log format to use",
				Validator: func(format string) error {
					formats := []string{"text", "json"}
					if !slices.Contains(formats, format) {
						return fmt.Errorf("unknown log format: %s", format)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "trace",
				Value: false,
				Usage: "enable OpenTelemetry tracing",
			},
		},
		Action: run,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	signal.Notify(sigHangup, syscall.SIGHUP)
	defer cancel()
	err := cmd.Run(ctx, os.Args)

	if err != nil {
		panic(err)
	}

}

type showbridgeApp struct {
	ctx          context.Context
	configPath   string
	logger       *slog.Logger
	router       *showbridge.Router
	routerRunner *sync.WaitGroup
	tracer       trace.Tracer
}

func readConfig(configPath string) (config.Config, error) {
	cfg := config.Config{}

	configBytes, err := os.ReadFile(configPath)

	if err != nil {
		return config.Config{}, err
	}

	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func run(ctx context.Context, c *cli.Command) error {
	configPath := c.String("config")
	if configPath == "" {
		return errors.New("config path cannot be empty")
	}

	logLevel := slog.LevelInfo

	logLevelFromFlag := c.String("log-level")

	switch logLevelFromFlag {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logHandlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	logOutput := os.Stderr

	var logHandler slog.Handler

	logFormat := c.String("log-format")

	switch logFormat {
	case "json":
		logHandler = slog.NewJSONHandler(logOutput, logHandlerOptions)
	case "text":
		logHandler = slog.NewTextHandler(logOutput, logHandlerOptions)
	default:
		logHandler = slog.NewTextHandler(logOutput, logHandlerOptions)
	}

	slog.SetDefault(slog.New(logHandler))

	var tracer trace.Tracer
	if c.Bool("trace") {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			return fmt.Errorf("failed to create trace exporter: %w", err)
		}

		tracerProvider := newTracerProvider(exporter)
		otel.SetTracerProvider(tracerProvider)
		defer tracerProvider.Shutdown(ctx)

		tracer = tracerProvider.Tracer("showbridge")
	} else {
		tracer = otel.Tracer("showbridge")
	}

	showbridgeApp := &showbridgeApp{
		ctx:          ctx,
		configPath:   configPath,
		logger:       slog.Default().With("component", "cmd"),
		routerRunner: &sync.WaitGroup{},
		tracer:       tracer,
	}

	router, err := showbridgeApp.getNewRouter()
	if err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}
	showbridgeApp.router = router

	showbridgeApp.routerRunner.Go(func() {
		router.Start(context.Background())
	})

	go showbridgeApp.handleHangup()

	<-showbridgeApp.ctx.Done()
	showbridgeApp.logger.Debug("shutting down router")
	showbridgeApp.router.Stop()
	showbridgeApp.logger.Debug("waiting for router to exit")
	showbridgeApp.routerRunner.Wait()
	return nil
}

func (app *showbridgeApp) handleHangup() {
	for {
		select {
		case <-sigHangup:
			app.logger.Info("received SIGHUP, reloading configuration")
			newRouter, err := app.getNewRouter()
			if err != nil {
				app.logger.Error("failed to reload configuration", "error", err)
				continue
			}
			app.router.Stop()
			app.routerRunner.Wait()
			app.router = newRouter
			app.routerRunner.Go(func() {
				app.router.Start(context.Background())
			})
			app.logger.Info("configuration reloaded successfully")
		case <-app.ctx.Done():
			return
		}
	}
}

func (app *showbridgeApp) getNewRouter() (*showbridge.Router, error) {
	// TODO(jwetzell): what should happen when the config file is unchanged?
	config, err := readConfig(app.configPath)
	if err != nil {
		return nil, err
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(config, app.tracer)

	for _, moduleError := range moduleErrors {
		app.logger.Error("problem initializing module", "index", moduleError.Index, "error", moduleError.Error)
	}

	for _, routeError := range routeErrors {
		app.logger.Error("problem initializing route", "index", routeError.Index, "error", routeError.Error)
	}

	if moduleErrors != nil || routeErrors != nil {
		return nil, fmt.Errorf("errors initializing modules or routes")
	}

	return router, nil
}

func newTracerProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("showbridge"),
			semconv.ServiceVersion(version),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
