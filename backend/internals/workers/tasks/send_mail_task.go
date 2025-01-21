package tasks

import (
	"context"
	"email-marketing-service/api/v1/utils"
	"email-marketing-service/internals/workers/payloads"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

// EnqueueEmail enqueues an email to be sent asynchronously
func EnqueueEmail(client *asynq.Client, subject, email, message, sender string) error {
	payload, err := json.Marshal(payloads.EmailPayload{
		Subject: subject,
		Email:   email,
		Message: message,
		Sender:  sender,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %v", err)
	}

	task := asynq.NewTask(TaskSendEmail, payload)
	_, err = client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Queue("emails"),
		asynq.Timeout(30*time.Second),
	)
	return err
}

// ProcessEmailTask handles the actual email sending
func ProcessEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload payloads.EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal email payload: %v", err)
	}

	// Get default SMTP config
	smtpConfig := utils.DefaultSMTPConfig()

	// Send the email using your existing sendMail function
	if err := utils.SendMail(
		payload.Subject,
		payload.Email,
		payload.Message,
		payload.Sender,
		&smtpConfig,
	); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
