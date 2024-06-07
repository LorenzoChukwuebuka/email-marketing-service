package main

import (
	"context"
	"email-marketing-service/api/database"
	"email-marketing-service/api/observers"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/routes"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var (
	response = &utils.ApiResponse{}
	//logger   = &utils.Logger{}
)

func cronJobs(dbConn *gorm.DB) *cron.Cron {
	subscriptionRepo := repository.NewSubscriptionRepository(dbConn)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	// Create a new cron scheduler
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		subscriptionService.UpdateExpiredSubscription()
	})

	return c
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
				// Log the panic
				fmt.Println("Recovered from panic:", r)
				// Print the stack trace
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				fmt.Printf("Panic Stack Trace:\n%s\n", stack)
				// Respond with an internal server error

				errorStack := map[string]interface{}{
					"Message": "Internal Server Error",
				}

				response.ErrorResponse(w, errorStack)

				//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func main() {

	logger, err := utils.NewLogger("app.log")

	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Initialize the database connection
	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	smtpWebHookRepo := repository.NewSMTPWebHookRepository(dbConn)

	eventBus := utils.GetEventBus()

	dbObserver := observers.NewCreateEmailStatusObserver(smtpWebHookRepo)
	eventBus.Register("send_success", dbObserver)
	eventBus.Register("send_failed", dbObserver)

	//instantiate the cron scheduler
	c := cronJobs(dbConn)

	r := mux.NewRouter()

	// Create a subrouter with the "/api/v1" prefix
	apiV1Router := r.PathPrefix("/api/v1").Subrouter()
	adminRouter := r.PathPrefix("/api/v1/admin").Subrouter()
	apiV1Router.Use(enableCORS)
	adminRouter.Use(enableCORS)
	routes.RegisterUserRoutes(apiV1Router, dbConn)
	routes.RegisterAdminRoutes(adminRouter, dbConn)

	r.Use(recoveryMiddleware)

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}

	// Start the cron scheduler in a goroutine
	go func() {
		// Start the cron scheduler
		c.Start()
		defer c.Stop()
		// Let it run indefinitely
		select {}
	}()

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
