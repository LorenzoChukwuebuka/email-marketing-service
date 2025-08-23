package main

import (
	"context"
	"database/sql"
	"email-marketing-service/core/server/asynqserver"
	"email-marketing-service/core/server/httpserver"
	"email-marketing-service/core/server/websocketserver"
	"email-marketing-service/internal/config"
	seeders "email-marketing-service/internal/db/seeder"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/logger"
	smtp_server "email-marketing-service/internal/smtp"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	// Run seeders
	if err := runSeeders(ctx, store); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start()
	}()

	//Start the SMTP server with functional options
	wg.Add(1)
	go func() {
		defer wg.Done()

		mode := os.Getenv("SERVER_MODE")
		var err error

		switch mode {
		case "production":
			err = smtp_server.StartSMTPServer(ctx, store,
				smtp_server.WithPort("587"),
				smtp_server.WithTLS(
					"/etc/letsencrypt/live/smtp.crabmailer.com/fullchain.pem",
					"/etc/letsencrypt/live/smtp.crabmailer.com/privkey.pem",
				),
				smtp_server.WithInsecureAuth(true),
			)
		case "staging":
			err = smtp_server.StartSMTPServer(ctx, store,
				smtp_server.WithPort("587"),
				smtp_server.WithTLS(
					"/etc/letsencrypt/live/smtp.crabmailer.com/fullchain.pem",
					"/etc/letsencrypt/live/smtp.crabmailer.com/privkey.pem",
				),
				smtp_server.WithInsecureAuth(true),
			)

		default: // development
			err = smtp_server.StartSMTPServer(ctx, store,
				smtp_server.WithInsecureAuth(true),
				// Port will use default from config
			)
		}

		if err != nil {
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

func runSeeders(ctx context.Context, store db.Store) error {
	// Create seeder manager
	manager := seeders.NewManager()

	// Register all seeders
	manager.RegisterAll(seeders.GetAllSeeders()...)

	// Run all seeders
	return manager.SeedAll(ctx, store)
}
