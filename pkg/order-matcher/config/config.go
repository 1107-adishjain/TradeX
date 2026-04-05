package config


import(
	"github.com/godotenv/godotenv"
	"os"
)

type Config struct {
	OrderMatcherPort string
}

func LoadConfig() (*Config, error){
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		OrderMatcherPort: getEnv("ORDER_MATCHER_PORT", "8083"),
	}, nil
}


func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}