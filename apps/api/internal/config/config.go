package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DatabaseUrl		string		`env:"DATABASE_URL"`
}

func LoadConfigs() (DatabaseConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return DatabaseConfig{}, err
	}

	var config DatabaseConfig
	if err := env.Parse(&config); err != nil {
		return DatabaseConfig{}, err
	}

	return config, nil
}
