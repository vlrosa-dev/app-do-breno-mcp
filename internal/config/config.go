package config

import "os"

// Defaults embutidos no binário. Podem ser sobrescritos por variável de ambiente.
const (
	DefaultAPIBaseURL   = "https://biijnyhamoyzanriocil.supabase.co"
	DefaultPublicAPIKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImJpaWpueWhhbW95emFucmlvY2lsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTQ4Njg3MzMsImV4cCI6MjA3MDQ0NDczM30.nlgHA3vqlehWSdNQmYCBmA33VLD9TCNYYenyiRtzsLg"
)

// Config reúne os valores efetivos usados para falar com a API.
type Config struct {
	APIBaseURL string
	APIKey     string
	Email      string
	Password   string
}

// Load resolve a config a partir do ambiente, caindo nos defaults embutidos
// quando a variável não está definida.
func Load() Config {
	return Config{
		APIBaseURL: firstNonEmpty(os.Getenv("API_BASE_URL"), DefaultAPIBaseURL),
		APIKey:     firstNonEmpty(os.Getenv("API_KEY"), DefaultPublicAPIKey),
		Email:      os.Getenv("API_USERNAME"),
		Password:   os.Getenv("API_PASSWORD"),
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
