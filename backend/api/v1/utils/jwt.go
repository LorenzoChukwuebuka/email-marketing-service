package utils

import (
	 
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt"
)

var (
	config = LoadEnv()
	key    = config.JWTKey
)

// Access token expiration (e.g., 15 minutes)
const accessTokenExpiration = time.Hour * 15 // 15 minutes

// Refresh token expiration (e.g., 7 days)
const refreshTokenExpiration = time.Hour * 24 * 7 // 7 days

// func JWTEncode(userId string, user_uuid string, username string, email string) (string, error) {

// 	// Create a new token object with claims
// 	claims := jwt.MapClaims{
// 		"sub":      "The server",
// 		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
// 		"username": username,                              // Include username claim
// 		"email":    email,
// 		"uuid":     user_uuid,
// 		"userId":   userId, // Include userId claim

// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Sign the token with a secret key

// 	tokenString, err := token.SignedString([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }


// GenerateAccessToken generates an access token with a short expiration time
func GenerateAccessToken(userId string, user_uuid string, username string, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      "The server",
		"exp":      time.Now().Add(accessTokenExpiration).Unix(), // 15 minutes expiration
		"username": username,
		"email":    email,
		"uuid":     user_uuid,
		"userId":   userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token with a long expiration time
func GenerateRefreshToken(userId string, user_uuid string, username string, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      "The server",
		"exp":      time.Now().Add(refreshTokenExpiration).Unix(), // 7 days expiration
		"username": username,
		"email":    email,
		"uuid":     user_uuid,
		"userId":   userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// If the token is valid, return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}

}

func GenerateAdminAccessToken(userId string, user_uuid string, admintype string, email string) (string, error) {
	LoadEnv()
	// Create a new token object with claims
	claims := jwt.MapClaims{
		"sub":    "The server",
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours                            // Include username claim
		"email":  email,
		"uuid":   user_uuid,
		"userId": userId, // Include userId claim
		"type":   admintype,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token with a long expiration time
func GenerateAdminRefreshToken(userId string, user_uuid string, admintype string, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":    "The server",
		"exp":    time.Now().Add(refreshTokenExpiration).Unix(), // 7 days expiration
		"type":   admintype,
		"email":  email,
		"uuid":   user_uuid,
		"userId": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
