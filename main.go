package main

import (
	"email-marketing-service/api/database"
	"email-marketing-service/api/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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
	routes.RegisterRoutes(apiV1Router)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe("localhost:9000", r))
}
