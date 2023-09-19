package utils

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

var key = os.Getenv("JWT_KEY")

func JWTEncode(userId int, user_uuid string, username string, email string) (string, error) {
	LoadEnv()
	// Create a new token object with claims
	claims := jwt.MapClaims{
		"sub":      "The server",
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"username": username,                              // Include username claim
		"email":    email,
		"uuid":     user_uuid,
		"userId":   userId, // Include userId claim

	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AdminJWTEncode(userId int, user_uuid string, admintype string, email string) (string, error) {
	LoadEnv()
	// Create a new token object with claims
	claims := jwt.MapClaims{
		"sub":      "The server",
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours                            // Include username claim
		"email":    email,
		"uuid":     user_uuid,
		"userId":   userId, // Include userId claim
		"type":     admintype,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return ""
	}

	return tokenParts[1]
}
