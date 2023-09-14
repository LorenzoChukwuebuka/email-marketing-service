package main

import (
	"context"
	"email-marketing-service/api/database"
	"email-marketing-service/api/routes"
	"email-marketing-service/api/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	response = &utils.ApiResponse{}
)

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

	fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)

	res := map[string]string{
		"response": "404 Not Found",
		"path":     r.URL.Path,
	}

	response.ErrorResponse(w, res)
}

func main() {
	// Initialize the database connection
	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbConn.Close()

	r := mux.NewRouter()

	// Create a subrouter with the "/api/v1" prefix
	apiV1Router := r.PathPrefix("/api/v1").Subrouter()
	adminRouter := r.PathPrefix("/api/v1/admin").Subrouter()
	apiV1Router.Use(enableCORS)
	adminRouter.Use(enableCORS)
	routes.RegisterUserRoutes(apiV1Router, dbConn)
	routes.RegisterAdminRoutes(adminRouter, dbConn)

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
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

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	fmt.Println("Server shut down gracefully")
}
