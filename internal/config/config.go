package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseUrl string
	SecreteKey  string
}

func MustLoad() *Config {

	// load env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env file required")
	}

	// check for port
	port := os.Getenv("port")
	if port == "" {
		panic("port is requuired")
	}

	// check for env
	env := os.Getenv("env")
	if env == "" {
		panic("env is required")
	}

	// check database url
	dbUrl := os.Getenv("databaseUrl")
	if dbUrl == "" {
		panic("database is required")
	}

	// cehck screte key
	key := os.Getenv("secreteString")
	if key == "" {
		panic("secrete key is required")
	}

	return &Config{
		Port:        port,
		Env:         env,
		DatabaseUrl: dbUrl,
		SecreteKey:  key,
	}
}
