package main

import (
	"context"
	"encoding/json"
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
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"
	"github.com/jwetzell/showbridge-go/internal/schema"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
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
	routerMutex  sync.Mutex
}

func readConfig(configPath string) (config.Config, error) {
	cfg := config.Config{}

	configBytes, err := os.ReadFile(configPath)

	if err != nil {
		return config.Config{}, err
	}

	//TODO(jwetzell): this is an annoying amount of marshaling

	yamlMap := make(map[string]any)

	err = yaml.Unmarshal(configBytes, &yamlMap)
	if err != nil {
		return config.Config{}, err
	}

	err = schema.ApplyDefaults(&yamlMap)
	if err != nil {
		return config.Config{}, err
	}

	err = schema.ValidateConfig(yamlMap)
	if err != nil {
		return config.Config{}, err
	}

	validatedConfigBytes, err := json.Marshal(yamlMap)

	err = json.Unmarshal(validatedConfigBytes, &cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func writeConfig(configPath string, newConfig config.Config) error {
	configBytes, err := yaml.Marshal(newConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, configBytes, 0644)
	if err != nil {
		return err
	}

	return nil
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

	if c.Bool("trace") {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			return fmt.Errorf("failed to create trace exporter: %w", err)
		}

		tracerProvider := newTracerProvider(exporter)
		otel.SetTracerProvider(tracerProvider)
		defer tracerProvider.Shutdown(ctx)

		otel.SetTracerProvider(tracerProvider)
	}

	showbridgeApp := &showbridgeApp{
		ctx:          ctx,
		configPath:   configPath,
		logger:       slog.Default().With("component", "cmd"),
		routerRunner: &sync.WaitGroup{},
	}

	config, err := readConfig(showbridgeApp.configPath)
	if err != nil {
		return err
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(config)

	showbridgeApp.logConfigErrors(moduleErrors, routeErrors)

	if moduleErrors != nil || routeErrors != nil {
		return fmt.Errorf("errors initializing modules or routes")
	}

	if err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}
	showbridgeApp.routerMutex.Lock()
	showbridgeApp.router = router

	showbridgeApp.routerRunner.Go(func() {
		router.Start(context.Background())
	})
	showbridgeApp.routerMutex.Unlock()

	go showbridgeApp.handleChannels()

	<-showbridgeApp.ctx.Done()
	showbridgeApp.logger.Debug("shutting down router")
	showbridgeApp.router.Stop()
	showbridgeApp.logger.Debug("waiting for router to exit")
	showbridgeApp.routerRunner.Wait()
	return nil
}

func (app *showbridgeApp) handleChannels() {
	for {
		select {
		case <-sigHangup:
			app.logger.Info("received SIGHUP, reloading configuration")
			app.routerMutex.Lock()
			config, err := readConfig(app.configPath)
			if err != nil {
				app.logger.Error("failed to read config file", "error", err)
				app.routerMutex.Unlock()
				continue
			}
			moduleErrors, routeErrors := app.router.UpdateConfig(config)
			app.logConfigErrors(moduleErrors, routeErrors)
			app.logger.Info("configuration reloaded successfully")
			app.routerMutex.Unlock()
		case config := <-app.router.ConfigChange:
			app.logger.Info("router config changed updating config file")
			err := writeConfig(app.configPath, config)
			if err != nil {
				app.logger.Error("failed to write config file", "error", err)
				continue
			}
			app.logger.Info("config file updated successfully")
		case <-app.ctx.Done():
			return
		}
	}
}

func (app *showbridgeApp) logConfigErrors(moduleErrors []module.ModuleError, routeErrors []route.RouteError) {
	for _, moduleError := range moduleErrors {
		app.logger.Error("problem initializing module", "index", moduleError.Index, "error", moduleError.Error)
	}

	for _, routeError := range routeErrors {
		app.logger.Error("problem initializing route", "index", routeError.Index, "error", routeError.Error)
	}
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
