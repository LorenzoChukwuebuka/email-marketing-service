package routes

import (
	"context"
	"email-marketing-service/api/controllers"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var key = os.Getenv("JWT_KEY")

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	utils.LoadEnv()
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

		if !ok {
			utils.ErrorResponse(w, "invalid jwt claims")
			return
		}

		ctx := context.WithValue(r.Context(), "jwtclaims", jwtclaims)
		// Proceed to the next handler
		next(w, r.WithContext(ctx))
	}
}

var (
	//intialize the dependencies
	otpRepo        = repository.NewOTPRepository()
	OTPService     = services.NewOTPService(otpRepo)
	UserRepo       = repository.NewUserRepository()
	UserServices   = services.NewUserService(UserRepo, OTPService)
	userController = controllers.NewUserController(UserServices)
)

var RegisterRoutes = func(router *mux.Router) {
	//userController := &controllers.UserController{}
	router.HandleFunc("/greet", JWTMiddleware(userController.Welcome)).Methods("GET")
	router.HandleFunc("/user-signup", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/verify-user", userController.VerifyUser).Methods("POST")
	router.HandleFunc("/user-login", userController.Login).Methods("POST")
	router.HandleFunc("/user-forget-password", userController.ForgetPassword).Methods("POST")
	router.HandleFunc("/user-reset-password", userController.ResetPassword).Methods("POST")
}
