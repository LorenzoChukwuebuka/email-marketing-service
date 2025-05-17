package main

import (
	"context"
	"database/sql"
	"email-marketing-service/core/server"
	"email-marketing-service/internal/config"
	seeders "email-marketing-service/internal/db/seeder"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var (
	cfg = config.LoadEnv()
)

func main() {

	conn, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing to database: %v", err)
		}
	}(conn)

	// Test the connection
	err = conn.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return
	}
	fmt.Println("Connected to the database!")

	store := db.NewStore(conn)

	server := server.NewServer(store)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seeders.SeedPlans(ctx, store)

	server.Start()

}
