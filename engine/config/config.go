package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUri  string
	ServerPort   string
	JWTSecret    string
	DB_SYNC      bool
	RABBITMQ_URL string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080" // default port if none is set
		}
	}

	return Config{
		DatabaseUri:  os.Getenv("DATABASE_URL"),
		ServerPort:   port,
		JWTSecret:    os.Getenv("JWT_SECRET"),
		DB_SYNC:      os.Getenv("DB_SYNC") == "true",
		RABBITMQ_URL: os.Getenv("RABBITMQ_URL"),
	}
}
