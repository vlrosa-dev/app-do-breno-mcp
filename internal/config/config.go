package config

import "os"

const (
	defaultAPIBaseURL   = "https://biijnyhamoyzanriocil.supabase.co"
	defaultPublicApiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImJpaWpueWhhbW95emFucmlvY2lsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTQ4Njg3MzMsImV4cCI6MjA3MDQ0NDczM30.nlgHA3vqlehWSdNQmYCBmA33VLD9TCNYYenyiRtzsLg"
)

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
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = defaultPublicApiKey
	}

	return &Config{
		APIBaseURL:  apiBaseURL,
		APIKey:      apiKey,
		APIUsername: os.Getenv("API_USERNAME"),
		APIPassword: os.Getenv("API_PASSWORD"),
	}
}
