package main

import (
	"context"
	"database/sql"
	"email-marketing-service/core/server/asynqserver"
	"email-marketing-service/core/server/httpserver"
	"email-marketing-service/core/server/websocketserver"
	"email-marketing-service/internal/logger"
	"email-marketing-service/internal/config"
	seeders "email-marketing-service/internal/db/seeder"
	db "email-marketing-service/internal/db/sqlc"
	smtp_server "email-marketing-service/internal/smtp"

	//"email-marketing-service/internal/workers/tasks"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := db.NewStore(conn)
	server := server.NewServer(store)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	//seeders
	seeders.SeedPlans(ctx, store)
	seeders.SeedSMTPKey(ctx, store)

	//redis...
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default for local development
	}

	asynqserver.Init()
	client := asynqserver.GetClient()
	defer client.Close()

	err = logger.InitDefaultLogger("logs/app.log", logger.INFO)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	//simulating mail sending
	//_ = tasks.EnqueueWelcomeEmail(client, "hello@hello.com", "obi")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start()
	}()

	//start the smtpserver
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := smtp_server.StartSMTPServer(ctx, store); err != nil {
			log.Printf("SMTP server error: %v", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := asynqserver.StartAsynqServer(redisAddr, store, ctx); err != nil {
			log.Printf("Asynq server error: %v", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		websocketserver.StartSocketServer()
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal")
	// Cancel the context to initiate shutdown of SMTP server
	cancel()
	// Wait for all components to shut down
	wg.Wait()
	log.Println("All components shut down gracefully")

}
