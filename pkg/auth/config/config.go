package config

import (
	"github.com/godotenv/godotenv"
	"os"
)

type Config struct {
	AuthPort string
	MongoURI string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		AuthPort: getEnv("AUTH_PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
