package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DATABASE_URL string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on OS environment variables")
	}

	return Config{
		Port:         os.Getenv("PORT"),
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	}
}
