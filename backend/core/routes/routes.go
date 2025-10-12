package routes

import (
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	worker "email-marketing-service/internal/workers"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func InitRoutes(r *mux.Router, store db.Store, wkr *worker.Worker) {
	// Apply middlewares
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.EnableCORS)
	r.Use(middleware.MethodNotAllowedMiddleware)
	r.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Health route
	apiV1.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Initialize all sub-routes
	routeMap := map[string]RouteInterface{
		"auth":      NewAuthRoute(store, wkr),
		"admin":     NewAdminRoute(store),
		"contacts":  NewContactRoutes(store),
		"templates": NewTemplateRoute(store),
		"payments":  NewPaymentRoute(store),
		"campaigns": NewCampaignRoute(store,wkr),
		"domains":   NewDomainRoute(store),
		"senders":   NewSenderRoute(store),
		"users":     NewUserRoute(store),
		"misc":      NewMiscRoute(store),
		"key":       NewSMTPAPIKeyRoute(store),
		"support":   NewSupportRoute(store),
	}

	for path, route := range routeMap {
		route.InitRoutes(apiV1.PathPrefix("/" + path).Subrouter())
	}

	// Serve uploads
	uploadsDir := filepath.Join(".", "uploads", "tickets")
	r.PathPrefix("/uploads/tickets/").Handler(
		http.StripPrefix("/uploads/tickets/", http.FileServer(http.Dir(uploadsDir))),
	)
}
