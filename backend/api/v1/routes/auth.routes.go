package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AuthRoute struct {
	db *gorm.DB
}

func NewAuthRoute(db *gorm.DB) *AuthRoute {
	return &AuthRoute{db: db}
}

func (ur *AuthRoute) InitRoutes(router *mux.Router) {
	userController, _ := InitializeUserController(ur.db)
	planController, _ := InitializePlanController(ur.db)
	sessionController, _ := InitializeUserssionController(ur.db)

	smtpController, _ := InitializeSMTPController(ur.db)
	transactionController, _ := InitializeTransactionController(ur.db)
	supportTicketController, _ := InitializeSupportTicketController(ur.db)
	subscriptionController, _ := InitializeSubscriptionController(ur.db)

	// auth routes
	router.HandleFunc("/greet", middleware.JWTMiddleware(userController.Welcome)).Methods("GET")
	router.HandleFunc("/user-signup", userController.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify-user", userController.VerifyUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/user-login", userController.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/user-forget-password", userController.ForgetPassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/user-reset-password", userController.ResetPassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/change-user-password", middleware.JWTMiddleware(userController.ChangeUserPassword)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/resend-otp", userController.ResendOTP).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-user-details", middleware.JWTMiddleware(userController.GetUserDetails)).Methods("GET", "OPTIONS")
	router.HandleFunc("/edit-user-details", middleware.JWTMiddleware(userController.EditUser)).Methods("PUT", "OPTIONS")

	// Transaction routes
	router.HandleFunc("/initialize-transaction", middleware.JWTMiddleware(transactionController.InitiateNewTransaction)).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify-transaction/{paymentmethod}/{reference}", middleware.JWTMiddleware(transactionController.ChargeTransaction)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-single-billing/{billingId}", middleware.JWTMiddleware(transactionController.GetSingleBillingRecord)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-all-billing", middleware.JWTMiddleware(transactionController.GetAllUserBilling))

	// Plan routes
	router.HandleFunc("/get-all-plans", planController.GetAllPlans).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-single-plan/{id}", planController.GetSinglePlan).Methods("GET", "OPTIONS")

	// Subscription routes
	router.HandleFunc("/cancel-subscription/{subscriptionId}", middleware.JWTMiddleware(subscriptionController.CancelSubscription)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-user-current-sub", middleware.JWTMiddleware(userController.GetUserSubscription)).Methods("GET", "OPTIONS")

	// SMTP routes
	router.HandleFunc("/smtp/email", smtpController.SendSMTPMail).Methods("POST", "OPTIONS")

	// Session routes
	router.HandleFunc("/create-session", sessionController.CreateSessions).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-sessions", middleware.JWTMiddleware(sessionController.GetAllSessions)).Methods("GET", "OPTIONS")
	router.HandleFunc("/delete-session", middleware.JWTMiddleware(sessionController.DeleteSession)).Methods("DELETE", "OPTIONS")

	// Support ticket routes
	router.HandleFunc("/create-ticket", middleware.JWTMiddleware(supportTicketController.CreateTicket)).Methods("POST", "OPTIONS")

	// Testing API routes
	router.HandleFunc("/update-expired-subscriptions", subscriptionController.UpdateAllExpiredSubscriptions).Methods("GET", "OPTIONS")
	router.HandleFunc("/test-create-daily-mail-calc", smtpController.CreateRecordDailyMailCalculation).Methods("POST", "OPTIONS")

}
