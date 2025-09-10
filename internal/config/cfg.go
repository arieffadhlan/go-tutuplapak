package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Database Database
}

func LoadsAllAppConfig() (*Config, error) {
	_ = godotenv.Load()
	
	return &Config{
		App:      loadAppConfig(),
		Database: loadDatabaseConfig(),
	}, nil
}
