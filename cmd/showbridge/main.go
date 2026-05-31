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
	"github.com/jwetzell/showbridge-go/internal/schema"
	"github.com/urfave/cli/v3"
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
				Name:    "config",
				Value:   "./config.yaml",
				Usage:   "path to config file",
				Sources: cli.EnvVars("SHOWBRIDGE_CONFIG"),
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
				Sources: cli.EnvVars("SHOWBRIDGE_LOG_LEVEL"),
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
				Sources: cli.EnvVars("SHOWBRIDGE_LOG_FORMAT"),
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
	ctx         context.Context
	configPath  string
	logger      *slog.Logger
	router      *showbridge.Router
	routerMutex sync.Mutex
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
		return config.Config{}, fmt.Errorf("failed to apply defaults: %w", err)
	}

	err = schema.ValidateConfig(yamlMap)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to validate config: %w", err)
	}

	validatedConfigBytes, err := json.Marshal(yamlMap)
	if err != nil {
		return config.Config{}, err
	}

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

	err = os.WriteFile(configPath, configBytes, 0600)
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

	var logLevel slog.Level

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

	showbridgeApp := &showbridgeApp{
		ctx:        ctx,
		configPath: configPath,
		logger:     slog.Default().With("component", "cmd"),
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

	showbridgeApp.routerMutex.Lock()
	showbridgeApp.router = router

	router.Start(context.Background())
	showbridgeApp.routerMutex.Unlock()

	go showbridgeApp.handleChannels()

	<-showbridgeApp.ctx.Done()
	showbridgeApp.logger.Debug("shutting down router")
	showbridgeApp.router.Stop()
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
			moduleErrors, routeErrors, err := app.router.UpdateConfig(config, false)
			if err != nil {
				app.logger.Error("failed to update router config", "error", err)
				app.routerMutex.Unlock()
				continue
			}
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

func (app *showbridgeApp) logConfigErrors(moduleErrors []config.ModuleError, routeErrors []config.RouteError) {
	for _, moduleError := range moduleErrors {
		app.logger.Error("problem initializing module", "index", moduleError.Index, "error", moduleError.Error)
	}

	for _, routeError := range routeErrors {
		app.logger.Error("problem initializing route", "index", routeError.Index, "error", routeError.Error)
	}
}
