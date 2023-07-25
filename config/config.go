package config

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
)
type AppConfig struct {
	Session *session.Store
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Could not load environment")
	}
}
