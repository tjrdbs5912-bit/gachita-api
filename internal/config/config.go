package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration loaded from the environment.
type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
