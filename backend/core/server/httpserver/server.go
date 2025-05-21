package server

import (
	"context"
	"email-marketing-service/core/middleware"
	"email-marketing-service/core/routes"
	"email-marketing-service/internal/config"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Server struct {
	router *mux.Router
	db     db.Store
}

func NewServer(db db.Store) *Server {
	return &Server{
		router: mux.NewRouter(),
		db:     db,
	}
}

var (
	cfg = config.LoadEnv()
)

func (s *Server) setupLogger() (*os.File, error) {
	logFile, err := os.OpenFile("logs/requests.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
		return nil, err
	}

	logger := log.New(logFile, "", log.LstdFlags)
	s.router.Use(middleware.LoggingMiddleware(logger))
	return logFile, nil
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.RecoveryMiddleware)
	s.router.Use(middleware.EnableCORS)
	s.router.Use(middleware.MethodNotAllowedMiddleware)
	s.router.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)
	apiV1 := s.router.PathPrefix("/api/v1").Subrouter()

	// Health route
	apiV1.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	routeMap := map[string]routes.RouteInterface{
		"auth":            routes.NewAuthRoute(s.db),
		"admin/auth":      routes.NewAdminAuthRoute(s.db),
		"admin/systemdns": routes.NewAdminDNSRoute(s.db),
		"contacts":        routes.NewContactRoutes(s.db),
	}

	for path, route := range routeMap {
		route.InitRoutes(apiV1.PathPrefix("/" + path).Subrouter())
	}

	uploadsDir := filepath.Join(".", "uploads", "tickets")
	s.router.PathPrefix("/uploads/tickets/").Handler(http.StripPrefix("/uploads/tickets/", http.FileServer(http.Dir(uploadsDir))))

	// mode := os.Getenv("SERVER_MODE")

	// var staticDir string

	// if mode == "" {
	// 	staticDir = "./client"
	// } else if mode == "test" {
	// 	staticDir = "/app/client"
	// } else {
	// 	staticDir = "/app/client"
	// }

	// // Handle static files using Gorilla Mux
	// s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// // Handle all other routes by serving index.html for the SPA
	// s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	path := filepath.Join(staticDir, r.URL.Path)

	// 	// If the requested file exists, serve it
	// 	if _, err := os.Stat(path); err == nil {
	// 		http.ServeFile(w, r, path)
	// 		return
	// 	}

	// 	// Otherwise, serve index.html
	// 	http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	// })
}

func (s *Server) Start() {
	s.setupRoutes()
	logFile, err := s.setupLogger()
	if err != nil {
		log.Fatal("Failed to set up logger:", err)
	}
	defer logFile.Close()

	// c := cronjobs.SetupCronJobs(s.db)

	// go func() {
	// 	c.Start()
	// 	defer c.Stop()
	// 	select {}
	// }()

	server := &http.Server{
		Addr:         cfg.APP_PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}

	go func() {
		log.Println("Server started on port " + cfg.APP_PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	//graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	fmt.Printf("Received signal: %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Print("Server shut down gracefully")
}
