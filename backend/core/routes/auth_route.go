package routes

import (
	"email-marketing-service/core/handler/auth/controller"
	"email-marketing-service/core/handler/auth/service"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type AuthRoute struct {
	store db.Store
}

func NewAuthRoute(store db.Store) RouteInterface {
	return &AuthRoute{
		store: store,
	}
}
func (a *AuthRoute) InitRoutes(r *mux.Router) {
	authService := service.NewAuthService(a.store)
	authController := controller.NewAuthController(authService)
	r.HandleFunc("/testing", authController.Welcome).Methods("GET", "OPTIONS")
	r.HandleFunc("/signup", authController.SignUp).Methods("POST", "OPTIONS")
	r.HandleFunc("/verify", authController.VerifyEmail).Methods("POST", "OPTIONS")
	r.HandleFunc("/resend-verification", authController.ResendVerificationEmail).Methods("POST", "OPTIONS")
	r.HandleFunc("/forget-password", authController.ForgotPassword).Methods("POST", "OPTIONS")
	r.HandleFunc("/reset-password", authController.ResetPassword).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", authController.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/change-password", authController.ChangePassword).Methods("POST", "OPTIONS")
	r.HandleFunc("/refresh-token", authController.RefreshTokenHandler).Methods("POST", "OPTIONS")
}
