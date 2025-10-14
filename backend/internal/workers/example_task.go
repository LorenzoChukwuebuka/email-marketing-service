package worker

import (
	"context"
	"log"
	"time"
)

/**
##################################
This is an example and a dummy.
Only meant as an example.
##################################
**/

func ProcessWelcomeEmailTask(ctx context.Context, payload EmailPayload) error {
	log.Printf("Sending welcome email to %s (%s)", payload.Subject, payload.Email)

	// TODO: Implement actual email sending logic
	// Example:
	// return w.emailClient.Send(ctx, EmailRequest{
	//     To:      payload.Email,
	//     Subject: "Welcome to our platform!",
	//     Body:    fmt.Sprintf("Hello %s, welcome aboard!", payload.Name),
	// })

	// Simulate work
	time.Sleep(100 * time.Millisecond)
	return nil
}
