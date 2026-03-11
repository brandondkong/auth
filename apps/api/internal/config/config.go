package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)


var configs *Config

// Singleton
func LoadConfigs() (Config, error) {
	if configs != nil {
		return *configs, nil
	}

	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	configs = &config
	return config, nil
}
