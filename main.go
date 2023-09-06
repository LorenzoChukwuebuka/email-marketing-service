package main

import (
	"email-marketing-service/api/database"
	"email-marketing-service/api/routes"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
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

func main() {

	// Initialize the database connection
	dbConn, err := database.InitDB()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return
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
	http.Handle("/", r)

	// Define the port
	port := 9000

	// Start the server
	fmt.Printf("Server started on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
