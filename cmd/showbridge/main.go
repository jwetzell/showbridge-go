package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
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
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "set log level to DEBUG",
			},
			&cli.BoolFlag{
				Name:  "json",
				Value: false,
				Usage: "log using JSON",
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

	if c.Bool("debug") {
		logLevel = slog.LevelDebug
	}

	logHandlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	logOutput := os.Stderr

	var logHandler slog.Handler = slog.NewTextHandler(logOutput, logHandlerOptions)

	if c.Bool("json") {
		logHandler = slog.NewJSONHandler(logOutput, logHandlerOptions)
	}

	slog.SetDefault(slog.New(logHandler))

	commandLogger := slog.Default().With("component", "cmd")

	router, moduleErrors, routeErrors := showbridge.NewRouter(context.Background(), config)

	for _, moduleError := range moduleErrors {
		commandLogger.Error("problem initializing module", "index", moduleError.Index, "error", moduleError.Error)
	}

	for _, routeError := range routeErrors {
		commandLogger.Error("problem initializing route", "index", routeError.Index, "error", routeError.Error)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Run()
	})

	<-ctx.Done()
	commandLogger.Debug("shutting down router")
	router.Stop()
	commandLogger.Debug("waiting for router to exit")
	routerRunner.Wait()
	return nil
}
