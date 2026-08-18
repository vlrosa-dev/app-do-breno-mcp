package config

import "os"

type Config struct {
	APIBaseURL  string `env:"API_BASE_URL"`
	APIUsername string `env:"API_EMAIL"`
	APIPassword string `env:"API_PASSWORD"`
	APIKey      string `env:"API_KEY"`
}

func Load() *Config {
	return &Config{
		APIBaseURL:  os.Getenv("API_BASE_URL"),
		APIUsername: os.Getenv("API_USERNAME"),
		APIPassword: os.Getenv("API_PASSWORD"),
		APIKey:      os.Getenv("API_KEY"),
	}
}
