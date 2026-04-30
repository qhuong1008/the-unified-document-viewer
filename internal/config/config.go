package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
	Env  string
}

func LoadConfig() *Config {
	return &Config{
		Host: getEnv("HOST", "localhost"),
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
