package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Env(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Load .env Error", err)
	}
	return os.Getenv(key)
}
