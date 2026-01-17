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

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/urfave/cli/v3"
	"sigs.k8s.io/yaml"
)

var (
	version = "dev"
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
		},
		Action: run,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := cmd.Run(ctx, os.Args)

	if err != nil {
		panic(err)
	}

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
		return errors.New("config value cannot be empty")
	}

	config, err := readConfig(configPath)
	if err != nil {
		return err
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

	commandLogger := slog.Default().With("component", "cmd")

	router, moduleErrors, routeErrors := showbridge.NewRouter(config)

	for _, moduleError := range moduleErrors {
		commandLogger.Error("problem initializing module", "index", moduleError.Index, "error", moduleError.Error)
	}

	for _, routeError := range routeErrors {
		commandLogger.Error("problem initializing route", "index", routeError.Index, "error", routeError.Error)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Run(context.Background())
	})

	<-ctx.Done()
	commandLogger.Debug("shutting down router")
	router.Stop()
	commandLogger.Debug("waiting for router to exit")
	routerRunner.Wait()
	return nil
}
