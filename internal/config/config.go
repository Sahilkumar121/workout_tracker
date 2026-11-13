package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DbUrl string
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

	return Config{
		Port:  port,
		DbUrl: databaseUrl,
	}
}
