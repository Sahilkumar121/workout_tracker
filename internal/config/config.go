package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DbUrl     string
	SecretKey string
}

func MustLoad() Config {
	if err := godotenv.Load(".env"); err != nil {
		panic("File .env not found")
	}

	port := os.Getenv("ENV_PORT")
	if strings.TrimSpace(port) == "" {
		panic("port is required")
	}

	databaseUrl := os.Getenv("ENV_DATABASE_URL")
	if strings.TrimSpace(databaseUrl) == "" {
		panic("Database url is required")
	}

	secreteKey := os.Getenv("ENV_SECRETKEY")
	if strings.TrimSpace(secreteKey) == "" {
		panic("Secret key is required")
	}

	return Config{
		Port:  port,
		DbUrl: databaseUrl,
	}
}
