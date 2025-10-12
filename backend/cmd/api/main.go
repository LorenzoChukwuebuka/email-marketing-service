package main

import (
	"context"
	"database/sql"
	"email-marketing-service/core/server"
	"email-marketing-service/internal/config"
	seeders "email-marketing-service/internal/db/seeder"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/logger"
	smtp_server "email-marketing-service/internal/smtp"
	worker "email-marketing-service/internal/workers"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := db.NewStore(conn)
	// Configure worker
	cfg := worker.WorkerConfig{
		MaxRetries:      3,
		RetryDelay:      5 * time.Second,
		PollInterval:    1 * time.Second,
		StaleTaskWindow: 10 * time.Minute,
	}
	// Create worker
	w := worker.NewWorker(store, cfg)
	w.Start(ctx, 5)

	server := server.NewServer(store, w)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

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

	// Example: Enqueue some tasks
	// go func() {
	// 	// Enqueue immediate task
	// 	payload := worker.EmailPayload{
	// 		Sender:  "user@example.com",
	// 		Subject: "Welcome!",
	// 		Message: "Welcome to our service",
	// 	}
	// 	taskID, err := w.EnqueueTask(ctx, worker.TaskSendWelcomeEmail, payload)
	// 	if err != nil {
	// 		log.Printf("Failed to enqueue task: %v", err)
	// 	} else {
	// 		log.Printf("Task enqueued with ID: %d", taskID)
	// 	}

	// 	// Enqueue delayed task (send after 1 hour)
	// 	delayedTaskID, err := w.EnqueueTaskDelayed(ctx, worker.TaskSendWelcomeEmail, payload, 1*time.Hour)
	// 	if err != nil {
	// 		log.Printf("Failed to enqueue delayed task: %v", err)
	// 	} else {
	// 		log.Printf("Delayed task enqueued with ID: %d", delayedTaskID)
	// 	}
	// }()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats, err := w.GetStats(ctx)
				if err != nil {
					log.Printf("Failed to get stats: %v", err)
					continue
				}
				slog.Info("Worker stats", slog.Any("worker stats", stats))
			}
		}
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
		server.StartSocketServer()
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal")
	// Cancel the context to initiate shutdown of SMTP server
	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := w.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
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
