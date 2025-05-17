package tasks

// import (
// 	"context"
// 	"email-marketing-service/api/v1/repository"
// 	adminrepository "email-marketing-service/api/v1/repository/admin"
// 	"email-marketing-service/api/v1/utils"
// 	"email-marketing-service/internals/workers/payloads"
// 	"encoding/json"
// 	"fmt"
// 	"github.com/hibiken/asynq"
// )

// func EnqueueUserNotification(client *asynq.Client, userId string, notificationTitle string) error {
// 	payload, err := json.Marshal(payloads.UserNotificationPayload{
// 		UserId:           userId,
// 		NotifcationTitle: notificationTitle,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal payload: %v", err)
// 	}

// 	task := asynq.NewTask(TaskSendUserNotification, payload)
// 	info, err := client.Enqueue(task)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
// 	return nil
// }

// func ProcessUserNotification(ctx context.Context, t *asynq.Task, userRepo *repository.UserNotificationRepository) error {
// 	fmt.Printf("Processing user notification task: id=%s\n", t.Type())

// 	var payload payloads.UserNotificationPayload
// 	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
// 		return fmt.Errorf("failed to unmarshal payload: %v", err)
// 	}
// 	fmt.Printf("Processing notification for user %s: %s\n", payload.UserId, payload.NotifcationTitle)

// 	if err := utils.CreateNotification(userRepo, payload.UserId, payload.NotifcationTitle); err != nil {
// 		return fmt.Errorf("failed to create notification: %v", err)
// 	}

// 	fmt.Println("Successfully processed user notification")
// 	return nil
// }

// func EnqueueAdminNotification(client *asynq.Client, userId string, notificationTitle string, link string) error {
// 	payload, err := json.Marshal(payloads.AdminNotificationPayload{UserId: userId, Link: link, NotificationTitle: notificationTitle})
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal payload: %v", err)
// 	}
// 	task := asynq.NewTask(TaskSendAdminNotification, payload)
// 	_, err = client.Enqueue(task)
// 	return err
// }

// func ProccessAdminNotification(ctx context.Context, t *asynq.Task, adminRepo *adminrepository.AdminNotificationRepository) error {
// 	var payload payloads.AdminNotificationPayload

// 	//unmarshal the payloads (more like the params)
// 	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
// 		return fmt.Errorf("failed to unmarshal payload: %v", err)
// 	}
// 	if err := utils.CreateAdminNotifications(adminRepo, payload.UserId, payload.Link, payload.NotificationTitle); err != nil {
// 		fmt.Printf("Failed to create notification: %v\n", err)
// 	}

// 	return nil
// }
