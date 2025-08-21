package asynqserver

import (
	"context"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/workers/tasks"
	"github.com/hibiken/asynq"
	"log"
)

func StartAsynqServer(redisAddr string, dbConn db.Store,ctx context.Context) error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 100, // Number of workers
		},
	)

	// Initialize repositories
	// userNotificationRepo := repository.NewUserNotificationRepository(dbConn)
	// adminNotificationRepo := adminrepository.NewAdminNoficationRepository(dbConn)

	// Create a ServeMux to register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TaskSendWelcomeEmail, tasks.ProcessWelcomeEmailTask)
	//mux.HandleFunc(TaskSendNotification, HandleNotification)

	// mux.HandleFunc(tasks.TaskSendUserNotification, func(ctx context.Context, t *asynq.Task) error {
	// 	return tasks.ProcessUserNotification(ctx, t, userNotificationRepo)
	// })

	// mux.HandleFunc(tasks.TaskSendAdminNotification, func(ctx context.Context, t *asynq.Task) error {
	// 	return tasks.ProccessAdminNotification(ctx, t, adminNotificationRepo)
	// })

	//_ = tasks.EnqueueWelcomeEmail(client, "hello@hello.com", "obi")

	// Start the server
	log.Println("Starting Asynq server...")
	if err := srv.Run(mux); err != nil {
		return err
	}
	return nil
}
