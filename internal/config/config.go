package config

import "os"

const defaultAPIBaseURL = "https://biijnyhamoyzanriocil.supabase.co"

type Config struct {
	APIBaseURL  string `env:"API_BASE_URL"`
	APIUsername string `env:"API_EMAIL"`
	APIPassword string `env:"API_PASSWORD"`
	APIKey      string `env:"API_KEY"`
}

func Load() *Config {
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = defaultAPIBaseURL
	}

	return &Config{
		APIBaseURL:  apiBaseURL,
		APIUsername: os.Getenv("API_USERNAME"),
		APIPassword: os.Getenv("API_PASSWORD"),
		APIKey:      os.Getenv("API_KEY"),
	}
}
