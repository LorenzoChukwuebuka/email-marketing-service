package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func JWTEncode(userId int, username string, email string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create a new token object with claims
	claims := jwt.MapClaims{
		"sub":      userId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"username": username,                              // Include username claim
		"email":    email,                                 // Include email claim
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	key := os.Getenv("JWT_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
