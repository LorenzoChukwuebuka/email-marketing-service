package main

import (
	"context"
	"email-marketing-service/api/v1"
	"email-marketing-service/api/v1/database"
	smtp_server "email-marketing-service/internals/smtp"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	server := v1.NewServer(dbConn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Start the API server
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start() // This function already includes graceful shutdown
	}()

	// Start the SMTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := smtp_server.StartSMTPServer(ctx, dbConn); err != nil {
			log.Printf("SMTP server error: %v", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()

		  v1.StartSocketServer() 
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
