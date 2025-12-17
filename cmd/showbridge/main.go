package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

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
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./config.yaml",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			configPath := c.String("config")
			if configPath == "" {
				return fmt.Errorf("config value cannot be empty")
			}

			config, err := readConfig(configPath)
			if err != nil {
				return err
			}
			router, moduleErrors, routeErrors := showbridge.NewRouter(ctx, config)
			for _, moduleError := range moduleErrors {
				slog.Error("problem initializing module", "index", moduleError.Index, "error", moduleError.Error)
			}

			for _, routeError := range routeErrors {
				slog.Error("problem initializing route", "index", routeError.Index, "error", routeError.Error)
			}
			router.Run()
			return nil
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Interrupt)
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
