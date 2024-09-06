package main

import (
	"context"
	"email-marketing-service/api/v1"
	"email-marketing-service/api/v1/database"
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/observers"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	smtp_server "email-marketing-service/smtp"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func cronJobs(dbConn *gorm.DB) *cron.Cron {
	subscriptionRepo := repository.NewSubscriptionRepository(dbConn)
	planRepo := repository.NewPlanRepository(dbConn)
	dailyRepo := repository.NewMailUsageRepository(dbConn)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, dailyRepo, planRepo)

	// Create a new cron scheduler
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		subscriptionService.UpdateExpiredSubscription()
	})

	return c
}

func main() {
	dto.InitValidate()

	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	smtpWebHookRepo := repository.NewMailStatusRepository(dbConn)
	eventBus := utils.GetEventBus()
	dbObserver := observers.NewCreateEmailStatusObserver(smtpWebHookRepo)
	eventBus.Register("send_success", dbObserver)
	eventBus.Register("send_failed", dbObserver)

	c := cronJobs(dbConn)

	go func() {
		c.Start()
		defer c.Stop()
		select {}
	}()

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

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal")

	// Cancel the context to initiate shutdown of SMTP server
	cancel()

	// Wait for all components to shut down
	wg.Wait()

	log.Println("All components shut down gracefully")
}
