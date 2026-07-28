package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr     string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()

	expiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_EXPIRES_IN: %w", err)
	}

	cfg := Config{
		HTTPAddr:     getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: expiresIn,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
