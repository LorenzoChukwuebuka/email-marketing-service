package cronjobs

import (
	"context"
	"fmt"
	"log"
	"time"
	db "email-marketing-service/internal/db/sqlc"
)

type AutoCloseSupportTicket struct {
	*BaseJob
}

func NewAutoCloseSupportTicket(store db.Store, ctx context.Context) *AutoCloseSupportTicket {
	baseJob := NewBaseJob(
		store,
		ctx,
		"auto_close_support_tickets",
		"AutoCloseSupportTicket",
		"Automatically close stale support tickets that have no replies for 48+ hours",
	)

	return &AutoCloseSupportTicket{
		BaseJob: baseJob,
	}
}

func (j *AutoCloseSupportTicket) Run() error {
	log.Println("Starting auto-close support ticket job...")

	// Close stale tickets (no reply for 48+ hours)
	closedTickets, err := j.store.CloseStaleTickets(j.ctx)
	if err != nil {
		return fmt.Errorf("error closing stale tickets: %w", err)
	}

	if len(closedTickets) == 0 {
		log.Println("No stale tickets found to close")
		return nil
	}

	log.Printf("Successfully auto-closed %d stale tickets:", len(closedTickets))

	// Log details of closed tickets and send notification emails
	emailFailures := 0
	for _, ticket := range closedTickets {
		var lastReplyInfo string
		if ticket.LastReply.Valid {
			lastReplyInfo = ticket.LastReply.Time.Format("2006-01-02 15:04:05")
		} else {
			lastReplyInfo = "never"
		}

		log.Printf("- Ticket #%s (ID: %s) - Last reply: %s",
			ticket.TicketNumber,
			ticket.ID.String(),
			lastReplyInfo,
		)

		// Send email notification to user about ticket closure
		err := j.sendTicketClosedEmail(ticket.Name, ticket.Email, ticket.TicketNumber, ticket.Subject)
		if err != nil {
			log.Printf("Failed to send closure email for ticket %s: %v", ticket.TicketNumber, err)
			emailFailures++
		}
	}

	log.Println("Auto-close support ticket job completed successfully")

	// Return an error if too many email failures occurred
	if emailFailures > len(closedTickets)/2 { // More than 50% failed
		return fmt.Errorf("high email failure rate: %d out of %d emails failed to send", emailFailures, len(closedTickets))
	}

	return nil
}

func (j *AutoCloseSupportTicket) Schedule() string {
	return "0 0 0 * * *" // Daily at midnight
}

// sendTicketClosedEmail sends an email notification when a ticket is auto-closed
func (j *AutoCloseSupportTicket) sendTicketClosedEmail(fullName, email, ticketNumber, subject string) error {
	// Email template for ticket closure
	emailSubject := fmt.Sprintf("Ticket #%s has been closed", ticketNumber)

	emailBody := fmt.Sprintf(`
		Dear %s,
		
		Your support ticket has been automatically closed due to inactivity.
		
		Ticket Details:
		- Ticket Number: %s
		- Subject: %s
		- Status: Closed
		- Closed Date: %s
		
		This ticket was closed because we haven't received a response from you in the past 48 hours.
		
		If you still need assistance with this issue, please feel free to create a new support ticket
		or reply to this email to reopen the ticket.
		
		Thank you for using our service.
		
		Best regards,
		Support Team
	`, fullName, ticketNumber, subject, time.Now().Format("January 2, 2006 at 3:04 PM"))

	fmt.Print(emailBody, emailSubject)

	// TODO: Replace with actual email service call
	// Example: return j.emailService.SendEmail(email, emailSubject, emailBody)

	// Simulate potential email failures for testing
	log.Printf("Would send email to %s for ticket %s", email, ticketNumber)
	return nil
}
