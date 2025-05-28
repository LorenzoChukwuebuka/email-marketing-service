package middleware

import (
	"context"
	"email-marketing-service/internal/config"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

var (
	key = config.LoadEnv().JWTKey
)

func AdminJWTMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := helper.ExtractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized: Token not found", http.StatusUnauthorized)
			return
		}

		// Define the secret key used for verification
		secretKey := []byte(key)

		// Parse and verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {

			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		jwtclaims, ok := token.Claims.(jwt.MapClaims)

		// for key, value := range jwtclaims {
		// 	fmt.Printf("%s: %v\n", key, value)
		// }

		if !ok {

			helper.ErrorResponse(w, fmt.Errorf("invalid jwt claims"),nil)
			return
		}

		claimType, ok := jwtclaims["type"].(string)
		if !ok || claimType != "admin" {

			http.Error(w, "Unauthorized: You are not an admin", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "adminclaims", jwtclaims)
		// Proceed to the next handler
		next(w, r.WithContext(ctx))
	}
}


func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := helper.ExtractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized: Token not found", http.StatusUnauthorized)
			return
		}

		secretKey := []byte(key)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Token not verified", http.StatusUnauthorized)
			return
		}

		jwtclaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helper.ErrorResponse(w, fmt.Errorf("invalid jwt claims"), nil)
			return
		}

		ctx := context.WithValue(r.Context(), "authclaims", jwtclaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := helper.ExtractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized: Token not found", http.StatusUnauthorized)
			return
		}

		// Define the secret key used for verification
		secretKey := []byte(key)

		// Parse and verify the token with MapClaims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Token not verified", http.StatusUnauthorized)
			return
		}

		// for key, value := range jwtclaims {
		// 	fmt.Printf("%s: %v\n", key, value)
		// }

		// Ensure claims are of type jwt.MapClaims
		jwtclaims, ok := token.Claims.(jwt.MapClaims)

		for key, value := range jwtclaims {
			fmt.Printf("%s: %v\n", key, value)
		}

		if !ok {
			helper.ErrorResponse(w, fmt.Errorf("invalid jwt claims"),nil)
			return
		}

		// Set the claims into the request context
		ctx := context.WithValue(r.Context(), "authclaims", jwtclaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
