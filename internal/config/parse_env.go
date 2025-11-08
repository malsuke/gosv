package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func (config *Config) ParseConfigFromEnv() error {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found", slog.String("error", fmt.Sprintf("%v", err)))
	}

	if err = env.Parse(config); err != nil {
		return err
	}

	return nil
}
