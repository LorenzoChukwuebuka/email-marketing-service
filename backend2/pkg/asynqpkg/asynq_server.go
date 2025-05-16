package asynqpkg

import (
	"context"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/internals/workers/tasks"
	"log"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func StartAsynqServer(redisAddr string, dbConn *gorm.DB) error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 100, // Number of workers
		},
	)

	// Initialize repositories
	userNotificationRepo := repository.NewUserNotificationRepository(dbConn)
	adminNotificationRepo := adminrepository.NewAdminNoficationRepository(dbConn)

	// Create a ServeMux to register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TaskSendWelcomeEmail, tasks.ProcessWelcomeEmailTask)
	//mux.HandleFunc(TaskSendNotification, HandleNotification)

	mux.HandleFunc(tasks.TaskSendUserNotification, func(ctx context.Context, t *asynq.Task) error {
		return tasks.ProcessUserNotification(ctx, t, userNotificationRepo)
	})

	mux.HandleFunc(tasks.TaskSendAdminNotification, func(ctx context.Context, t *asynq.Task) error {
		return tasks.ProccessAdminNotification(ctx, t, adminNotificationRepo)
	})

	// Start the server
	log.Println("Starting Asynq server...")
	if err := srv.Run(mux); err != nil {
		return err
	}
	return nil
}
