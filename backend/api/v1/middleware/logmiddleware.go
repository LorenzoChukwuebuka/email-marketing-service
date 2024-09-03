package middleware

import (
	 
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	 "github.com/golang-jwt/jwt"
	"net/http"
		"github.com/gorilla/mux"
)

 

func LoggingMiddleware(logService *repository.LogRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			tokenString := utils.ExtractTokenFromHeader(r)

			if tokenString != "" {
				// Define the secret key used for verification
				secretKey := []byte(key)

				// Parse and verify the token
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return secretKey, nil
				})
				if err == nil && token.Valid {
					// Extract the claims from the token
					jwtclaims, ok := token.Claims.(jwt.MapClaims)
					if ok {
						// Extract the user ID from claims
						userID, ok := jwtclaims["user_id"].(string)
						if ok {
							// Log the request
							err = logService.LogAction(userID, "HTTP "+r.Method, r.RequestURI)
							if err != nil {
								http.Error(w, "Could not log request", http.StatusInternalServerError)
								return
							}
						}
					}
				} else {
					http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
					return
				}
			} else {
				// Handle cases where there is no token (e.g., public endpoints)
				// Optional: Log a message or skip logging if no token is present
				fmt.Println("No token provided, skipping logging.")
			}

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
