package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB_User               string `env:"DB_USER"`
	DBPassword            string `env:"DB_PASSWORD"`
	DBName                string `env:"DB_NAME"`
	DBHost                string `env:"DB_HOST"`
	DBPort                string `env:"DB_PORT"`
	MailUsername          string `env:"MAIL_USERNAME"`
	MailPassword          string `env:"MAIL_PASSWORD"`
	JWTKey                string `env:"JWT_KEY"`
	PaystackKey           string `env:"PAYSTACK_KEY"`
	PaystackBaseURL       string `env:"PAYSTACK_BASE_URL"`
	DB_URL                string `env:"DB_URL"`
	APPName               string `env:"APP_NAME"`
	PAYSTACK_CALLBACK_URL string `env:"PAYSTACK_CALLBACK_URL"`
	ENC_KEY               string `env:"ENC_KEY"`
	SMTP_SERVER           string `env:"SMTP_SERVER"`
	SMTP_PORT             string `env:"SMTP_PORT"`
	MAIL_PROCESSOR        string `env:"MAIL_PROCESSOR"`
	ADMIN_EMAIL           string `env:"ADMIN_EMAIL"`
	ADMIN_PASSWORD        string `env:"ADMIN_PASSWORD"`
}

var LoadEnv = func() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		DB_User:               os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		MailUsername:          os.Getenv("MAIL_USERNAME"),
		MailPassword:          os.Getenv("MAIL_PASSWORD"),
		JWTKey:                os.Getenv("JWT_KEY"),
		PaystackKey:           os.Getenv("PAYSTACK_KEY"),
		PaystackBaseURL:       os.Getenv("PAYSTACK_BASE_URL"),
		DB_URL:                os.Getenv("DB_URL"),
		APPName:               os.Getenv("APP_NAME"),
		PAYSTACK_CALLBACK_URL: os.Getenv("PAYSTACK_CALLBACK_URL"),
		ENC_KEY:               os.Getenv("ENC_KEY"),
		SMTP_SERVER:           os.Getenv("SMTP_SERVER"),
		SMTP_PORT:             os.Getenv("SMTP_PORT"),
		MAIL_PROCESSOR:        os.Getenv("MAIL_PROCESSOR"),
		ADMIN_EMAIL:           os.Getenv("ADMIN_EMAIL"),
		ADMIN_PASSWORD:        os.Getenv("ADMIN_PASSWORD"),
	}

	return config
}
