package config

import (
	"github.com/godotenv/godotenv"
	"os"
	"strings"
)

type Config struct {
	MarketDataPort string
	KafkaBroker    string
	Symbols        []string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		MarketDataPort: getEnv("MARKET_DATA_PORT", "8082"),
		KafkaBroker:    getEnv("KAFKA_BROKER", "localhost:9092"),
		Symbols:        parseSymbols(getEnv("MARKET_SYMBOLS", "")),
	}, nil
}

func parseSymbols(raw string) []string {
	parts := strings.Split(raw, ",")
	symbols := make([]string, 0, len(parts))

	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s == "" {
			continue
		}
		symbols = append(symbols, strings.ToUpper(s))
	}

	if len(symbols) == 0 {
		return []string{"BTCUSDT"}
	}

	return symbols
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
