package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TransactionRoute struct {
	db *gorm.DB
}

func NewTransactionRoute(db *gorm.DB) *TransactionRoute {
	return &TransactionRoute{db: db}
}

func (ur *TransactionRoute) InitRoutes(router *mux.Router) {
	transactionController, _ := InitializeTransactionController(ur.db)
	router.HandleFunc("/initialize-transaction", middleware.JWTMiddleware(transactionController.InitiateNewTransaction)).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify-transaction/{paymentmethod}/{reference}", middleware.JWTMiddleware(transactionController.ChargeTransaction)).Methods("GET", "OPTIONS")
	//router.HandleFunc("/get-single-billing/{billingId}", middleware.JWTMiddleware(transactionController.GetSingleBillingRecord)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-all-billing", middleware.JWTMiddleware(transactionController.GetAllUserBilling))
}
