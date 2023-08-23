package utils

import (
	"github.com/joho/godotenv"
	"log"
)

var LoadEnv = func() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
