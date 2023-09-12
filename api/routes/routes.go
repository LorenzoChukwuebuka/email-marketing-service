package routes

import (
	"database/sql"
	"email-marketing-service/api/controllers"
	"email-marketing-service/api/middleware"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/services"
	"github.com/gorilla/mux"
)

var RegisterUserRoutes = func(router *mux.Router, db *sql.DB) {

	//intialize the user  dependencies
	otpRepo := repository.NewOTPRepository(db)
	OTPService := services.NewOTPService(otpRepo)
	UserRepo := repository.NewUserRepository(db)
	UserServices := services.NewUserService(UserRepo, OTPService)
	userController := controllers.NewUserController(UserServices)

	//plan
	planRepo := repository.NewPlanRepository(db)
	planService := services.NewPlanService(planRepo)
	planController := controllers.NewPlanController(planService)

	router.HandleFunc("/greet", middleware.JWTMiddleware(userController.Welcome)).Methods("GET")
	router.HandleFunc("/user-signup", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/verify-user", userController.VerifyUser).Methods("POST")
	router.HandleFunc("/user-login", userController.Login).Methods("POST")
	router.HandleFunc("/user-forget-password", userController.ForgetPassword).Methods("POST")
	router.HandleFunc("/user-reset-password", userController.ResetPassword).Methods("POST")
	router.HandleFunc("/change-user-password", middleware.JWTMiddleware(userController.ChangeUserPassword)).Methods("PUT")

	//public api
	router.HandleFunc("/get-all-plans", planController.GetAllPlans).Methods("GET")
	router.HandleFunc("/get-single-plan/{id}", planController.GetSinglePlan).Methods("GET")
}
