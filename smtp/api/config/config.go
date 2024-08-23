package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB_User    string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
}

var LoadEnv = func() *Config {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		DB_User:    os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}

	return config
}
