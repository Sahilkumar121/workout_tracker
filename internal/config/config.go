package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func MustLoad() Config {
	if err := godotenv.Load(".env"); err != nil {
		panic("File .env not found")
	}

	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		panic("port is required")
	}

	return Config{
		Port: port,
	}
}
