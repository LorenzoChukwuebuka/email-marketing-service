package main

import (
	"email-marketing-service/api/v1/database"
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/observers"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
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

func main() {
	logger, err := utils.NewLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

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

	server := NewServer(dbConn)
	server.Start()
}
