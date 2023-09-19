package middleware

import (
	"context"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

var key = os.Getenv("JWT_KEY")

func AdminJWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	utils.LoadEnv()
	response := &utils.ApiResponse{}
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.ExtractTokenFromHeader(r)
		if tokenString == "" {
			 
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	 
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		jwtclaims, ok := token.Claims.(jwt.MapClaims)

		// for key, value := range jwtclaims {
		// 	fmt.Printf("%s: %v\n", key, value)
		// }

		if !ok {

			response.ErrorResponse(w, "invalid jwt claims")
			return
		}

		claimType, ok := jwtclaims["type"].(string)
		if !ok || claimType != "admin" {
		 
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "adminclaims", jwtclaims)
		// Proceed to the next handler
		next(w, r.WithContext(ctx))
	}
}



func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	utils.LoadEnv()
	response := &utils.ApiResponse{}
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.ExtractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		jwtclaims, ok := token.Claims.(jwt.MapClaims)

		// for key, value := range jwtclaims {
		// 	fmt.Printf("%s: %v\n", key, value)
		// }

		if !ok {

			response.ErrorResponse(w, "invalid jwt claims")
			return
		}

		ctx := context.WithValue(r.Context(), "authclaims", jwtclaims)
		// Proceed to the next handler
		next(w, r.WithContext(ctx))
	}
}
