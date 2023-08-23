package utils

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

var key = os.Getenv("JWT_KEY")

func JWTEncode(userId int, username string, email string) (string, error) {
	LoadEnv()
	// Create a new token object with claims
	claims := jwt.MapClaims{
		"sub":      userId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"username": username,                              // Include username claim
		"email":    email,                                 // Include email claim
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
