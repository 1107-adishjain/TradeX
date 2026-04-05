package config


import(
	"github.com/godotenv/godotenv"
	"os"
)

type Config struct {
	NotifierPort string
}

func LoadConfig() (*Config, error){
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		NotifierPort: getEnv("NOTIFIER_PORT", "8081"),
	}, nil
}


func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}