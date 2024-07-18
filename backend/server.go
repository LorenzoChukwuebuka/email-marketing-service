package main

import (
	"context"
	"email-marketing-service/api/v1/routes"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Server struct {
	router *mux.Router
	db     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		router: mux.NewRouter(),
		db:     db,
	}
}

var (
	response = &utils.ApiResponse{}
	//logger   = &utils.Logger{}
)

func (s *Server) setupRoutes() {
	apiV1Router := s.router.PathPrefix("/api/v1").Subrouter()
	adminRouter := s.router.PathPrefix("/api/v1/admin").Subrouter()

	apiV1Router.Use(enableCORS)
	adminRouter.Use(enableCORS)

	routes.RegisterUserRoutes(apiV1Router, s.db)
	routes.RegisterAdminRoutes(adminRouter, s.db)

	s.router.Use(recoveryMiddleware)
	s.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}

func (s *Server) Start() {
	s.setupRoutes()

	server := &http.Server{
		Addr:    ":9000",
		Handler: s.router,
	}

	go func() {
		fmt.Println("Server started on port 9000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	fmt.Println("Server shut down gracefully")
}

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the request method is OPTIONS, just return a 200 status (pre-flight request)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the actual handler
		handler.ServeHTTP(w, r)
	})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	//fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)

	res := map[string]string{
		"response": "404 Not Found",
		"path":     r.URL.Path,
	}
	response.ErrorResponse(w, res)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				fmt.Printf("Panic Stack Trace:\n%s\n", stack)

				errorStack := map[string]interface{}{
					"Message": "Internal Server Error",
				}

				response.ErrorResponse(w, errorStack)

			}
		}()

		next.ServeHTTP(w, r)
	})
}
