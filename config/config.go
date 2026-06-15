package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	DBDriver string
	DBSource string
}

func Load() Config {
	return Config{
		Port:     getenv("PORT", ":3000"),
		DBDriver: getenv("DB_DRIVER", "postgres"),
		DBSource: getenv("DB_SOURCE", "postgres://postgres:postgres@localhost:5432/user_api?sslmode=disable"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c Config) Address() string {
	return fmt.Sprintf("%s", c.Port)
}
