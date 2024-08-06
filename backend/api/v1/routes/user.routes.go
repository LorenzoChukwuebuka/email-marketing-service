package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserRoute struct {
	db *gorm.DB
}

func NewUserRoute(db *gorm.DB) *UserRoute {
	return &UserRoute{db: db}
}

func (ur *UserRoute) InitRoutes(router *mux.Router) {
	userController, _ := InitializeUserController(ur.db)
	planController, _ := InitializePlanController(ur.db)
	sessionController, _ := InitializeUserssionController(ur.db)
	apiKeyController, _ := InitializeAPIKeyController(ur.db)
	smtpController, _ := InitializeSMTPController(ur.db)
	transactionController, _ := InitializeTransactionController(ur.db)
	supportTicketController, _ := InitializeSupportTicketController(ur.db)
	smptKeyController, _ := InitializeSMTPKeyController(ur.db)
	contactController, _ := InitializeContactController(ur.db)
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

	// API key routes
	router.HandleFunc("/generate-apikey", middleware.JWTMiddleware(apiKeyController.GenerateAPIKEY)).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete-apikey/{apiKeyId}", middleware.JWTMiddleware(apiKeyController.DeleteAPIKey)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/get-apikey", middleware.JWTMiddleware(apiKeyController.GetAPIKey)).Methods("GET", "OPTIONS")

	// SMTP key routes
	router.HandleFunc("/generate-new-smtp-master-password", middleware.JWTMiddleware(smptKeyController.GenerateNewSMTPMasterPassword)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-smtp-keys", middleware.JWTMiddleware(smptKeyController.GetUserSMTPKeys)).Methods("GET", "OPTIONS")
	router.HandleFunc("/create-smtp-key", middleware.JWTMiddleware(smptKeyController.CreateSMTPKey)).Methods("POST", "OPTIONS")
	router.HandleFunc("/toggle-smtp-key-status/{smtpKeyId}", middleware.JWTMiddleware(smptKeyController.ToggleSMTPKeyStatus)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/delete-smtp-key/{smtpKeyId}", middleware.JWTMiddleware(smptKeyController.DeleteSMTPKey)).Methods("DELETE", "OPTIONS")

	// Contact routes
	router.HandleFunc("/create-contact", middleware.JWTMiddleware(contactController.CreateContact)).Methods("POST", "OPTIONS")
	router.HandleFunc("/upload-contact-csv", middleware.JWTMiddleware(contactController.UploadContactViaCSV)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-all-contacts", middleware.JWTMiddleware(contactController.GetAllContacts))
	router.HandleFunc("/update-contact/{contactId}", middleware.JWTMiddleware(contactController.UpdateContact)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/delete-contact/{contactId}", middleware.JWTMiddleware(contactController.DeleteContact)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/create-contact-group", middleware.JWTMiddleware(contactController.CreateGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/add-contact-to-group", middleware.JWTMiddleware(contactController.AddContactToGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-all-contact-groups", middleware.JWTMiddleware(contactController.GetAllContactGroups))
	router.HandleFunc("/add-contact-to-group", middleware.JWTMiddleware(contactController.AddContactToGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/remove-contact-from-group", middleware.JWTMiddleware(contactController.RemoveContactFromGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-single-group/{groupId}", middleware.JWTMiddleware(contactController.GetASingleGroupWithContacts))
	router.HandleFunc("/delete-group/{groupId}", middleware.JWTMiddleware(contactController.DeleteContactGroup)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/edit-group/{groupId}", middleware.JWTMiddleware(contactController.UpdateContactGroup)).Methods("PUT", "OPTIONS")

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
