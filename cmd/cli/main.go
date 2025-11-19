package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jwetzell/showbridge-go"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "showbridge",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./config.json",
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
			router, err := showbridge.NewRouter(ctx, config)
			if err != nil {
				return err
			}
			router.Run()

			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)

	if err != nil {
		panic(err)
	}

}

func readConfig(configPath string) (showbridge.Config, error) {
	configBytes, err := os.ReadFile(configPath)

	if err != nil {
		return showbridge.Config{}, err
	}

	config := showbridge.Config{}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return showbridge.Config{}, err
	}

	return config, nil
}
