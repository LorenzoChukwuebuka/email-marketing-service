package server

import (
	"context"
	"email-marketing-service/core/middleware"
	"email-marketing-service/core/routes"
	"email-marketing-service/internal/config"
	"email-marketing-service/internal/cronjobs"
	db "email-marketing-service/internal/db/sqlc"
	worker "email-marketing-service/internal/workers"
	"fmt"
	"log"
	"net/http"
	 "log/slog"
	"os"
	"os/signal"
	"github.com/gorilla/mux"
	//"path/filepath"
	"syscall"
	"time"
)

type Server struct {
	router *mux.Router
	db     db.Store
	wkr    *worker.Worker
}

func NewServer(db db.Store, wkr *worker.Worker) *Server {
	return &Server{
		router: mux.NewRouter(),
		db:     db,
		wkr:    wkr,
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

func (s *Server) Start() {
	// Initialize routes
	routes.InitRoutes(s.router, s.db,s.wkr)
	logFile, err := s.setupLogger()
	if err != nil {
		log.Fatal("Failed to set up logger:", err)
	}
	defer logFile.Close()

	//start the cron jobs
	c := cronjobs.SetupCronJobs(s.db)
	c.Start()
	slog.Info("Cron jobs started")

	server := &http.Server{
		Addr:         cfg.APP_PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}

	go func() {
		slog.Info("Server started on port " + cfg.APP_PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	//graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	fmt.Printf("Received signal: %v\n", sig)

	// Stop cron jobs first
	slog.Info("Stopping cron jobs...")
	c.Stop()
	slog.Info("Cron jobs stopped")

	// Then shutdown the HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Print("Server shut down gracefully")
}
