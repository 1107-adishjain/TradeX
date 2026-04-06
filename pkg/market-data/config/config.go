package config

import (
	"github.com/godotenv/godotenv"
	"os"
)

type Config struct {
	MarketDataPort string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		MarketDataPort: getEnv("MARKET_DATA_PORT", "8082"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
