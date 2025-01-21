package tasks

import (
	"context"
	"email-marketing-service/internals/workers/payloads"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

func EnqueueWelcomeEmail(client *asynq.Client, email, name string) error {
	payload, err := json.Marshal(payloads.EmailExamplePayload{Email: email, Name: name})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}
	task := asynq.NewTask(TaskSendWelcomeEmail, payload)
	_, err = client.Enqueue(task)
	return err
}

func ProcessWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload payloads.EmailExamplePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	// Simulate sending email
	fmt.Printf("Sending welcome email to %s (%s)\n", payload.Name, payload.Email)
	return nil
}
