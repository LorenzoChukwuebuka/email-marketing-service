package routes

import (
	"email-marketing-service/core/handler/api_smtp_keys/controller"
	"email-marketing-service/core/handler/api_smtp_keys/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"

	"github.com/gorilla/mux"
)

type SMTPAPIKeyRoute struct {
	store db.Store
}

// Fixed: Return correct type
func NewSMTPAPIKeyRoute(store db.Store) *SMTPAPIKeyRoute {
	return &SMTPAPIKeyRoute{
		store: store,
	}
}

func (t *SMTPAPIKeyRoute) InitRoutes(r *mux.Router) {
	r.Use(middleware.JWTMiddleware)
	apiKeyRouter := r.PathPrefix("/apikey").Subrouter()
	apiKeysvc := services.NewAPIKeyService(t.store)
	apiKeyhandler := controller.NewAPIKeyController(apiKeysvc)
	{
		apiKeyRouter.HandleFunc("/create", apiKeyhandler.GenerateAPIKEY).Methods("POST", "OPTIONS")
		apiKeyRouter.HandleFunc("/get", apiKeyhandler.GetAPIKey).Methods("GET", "OPTIONS")
		apiKeyRouter.HandleFunc("/delete/{apikeyId}", apiKeyhandler.DeleteAPIKey).Methods("DELETE", "OPTIONS")
	}
	smtpkeyRouter := r.PathPrefix("/smtpkey").Subrouter()
	smtpkeysvc := services.NewSMTPKeyService(t.store)
	smtpkeyhandler := controller.NewSMTPKeyController(smtpkeysvc)
	{
		smtpkeyRouter.HandleFunc("/generate-new-masterkey", smtpkeyhandler.GenerateNewSMTPMasterPassword).Methods("POST", "OPTIONS")
		smtpkeyRouter.HandleFunc("/create", smtpkeyhandler.CreateSMTPKEY).Methods("POST", "OPTIONS")
		smtpkeyRouter.HandleFunc("/get", smtpkeyhandler.GetSMTPKey).Methods("GET", "OPTIONS")
		smtpkeyRouter.HandleFunc("/toggle-status/{smtpkeyId}", smtpkeyhandler.ToggleSMTPKeyStatus).Methods("PUT", "OPTIONS")
		smtpkeyRouter.HandleFunc("/delete/{smtpkeyId}", smtpkeyhandler.DeleteSMTPKey).Methods("DELETE", "OPTIONS")
	}
}
