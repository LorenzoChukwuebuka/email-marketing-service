package routes

import (
	"email-marketing-service/api/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var RegisterUserRoutes = func(router *mux.Router, db *gorm.DB) {
	userController, _ := InitializeUserController(db)
	planController, _ := InitializePlanController(db)
	sessionController, _ := InitializeUserssionController(db)
	apiKeyController, _ := InitializeAPIKeyController(db)
	smtpController, _ := InitializeSMTPController(db)
	transactionController, _ := InitializeTransactionController(db)
	supportTicketController, _ := InitializeSupportTicketController(db)

	//subscription service for testing only
	subscriptionController, _ := InitializeSubscriptionController(db)

	router.HandleFunc("/greet", middleware.JWTMiddleware(userController.Welcome)).Methods("GET")
	router.HandleFunc("/user-signup", userController.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify-user", userController.VerifyUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/user-login", userController.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/user-forget-password", userController.ForgetPassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/user-reset-password", userController.ResetPassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/change-user-password", middleware.JWTMiddleware(userController.ChangeUserPassword)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/resend-otp", userController.ResendOTP).Methods("POST", "OPTIONS")
	//transaction routes
	router.HandleFunc("/initialize-transaction", middleware.JWTMiddleware(transactionController.InitiateNewTransaction)).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify-transaction/{paymentmethod}/{reference}", middleware.JWTMiddleware(transactionController.ChargeTransaction)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-single-billing/{billingId}", middleware.JWTMiddleware(transactionController.GetSingleBillingRecord)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-all-billing", middleware.JWTMiddleware(transactionController.GetAllUserBilling)).Queries("page", "{page}").Methods("GET", "OPTIONS")

	//public api
	router.HandleFunc("/get-all-plans", planController.GetAllPlans).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-single-plan/{id}", planController.GetSinglePlan).Methods("GET", "OPTIONS")

	//subscription route
	router.HandleFunc("/cancel-subscription/{subscriptionId}", middleware.JWTMiddleware(subscriptionController.CancelSubscription)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-user-current-sub", middleware.JWTMiddleware(userController.GetUserSubscription)).Methods("GET", "OPTIONS")

	//api key route
	router.HandleFunc("/generate-apikey", middleware.JWTMiddleware(apiKeyController.GenerateAPIKEY)).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete-apikey/{apiKeyId}", middleware.JWTMiddleware(apiKeyController.DeleteAPIKey)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/get-apikey", middleware.JWTMiddleware(apiKeyController.GetAPIKey)).Methods("GET", "OPTIONS")

	//smtp
	router.HandleFunc("/smtp/email", smtpController.SendSMTPMail).Methods("POST", "OPTIONS")

	//session
	router.HandleFunc("/create-session", sessionController.CreateSessions).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-sessions", middleware.JWTMiddleware(sessionController.GetAllSessions)).Methods("GET", "OPTIONS")
	router.HandleFunc("/delete-session", middleware.JWTMiddleware(sessionController.DeleteSession)).Methods("DELETE", "OPTIONS")

	//support ticket

	router.HandleFunc("/create-ticket", middleware.JWTMiddleware(supportTicketController.CreateTicket)).Methods("POST", "OPTIONS")

	// Testing API
	router.HandleFunc("/update-expired-subscriptions", subscriptionController.UpdateAllExpiredSubscriptions).Methods("GET", "OPTIONS")
	router.HandleFunc("/test-create-daily-mail-calc", smtpController.CreateRecordDailyMailCalculation).Methods("POST", "OPTIONS")

}
