package cronjobs

import (
	"context"
	"log"
	db "email-marketing-service/internal/db/sqlc"
)

type AutoCloseSupportTicket struct {
	ctx   context.Context
	store db.Store
}

func NewAutoCloseSupportTicket(ctx context.Context, store db.Store) *AutoCloseSupportTicket {
	return &AutoCloseSupportTicket{
		ctx:   ctx,
		store: store,
	}
}

func (j *AutoCloseSupportTicket) Run() {
	log.Println("Starting auto-close support ticket job...")

	// Close stale tickets (no reply for 48+ hours)
	closedTickets, err := j.store.CloseStaleTickets(j.ctx)
	if err != nil {
		log.Printf("Error closing stale tickets: %v", err)
		return
	}

	if len(closedTickets) == 0 {
		log.Println("No stale tickets found to close")
		return
	}

	log.Printf("Successfully auto-closed %d stale tickets:", len(closedTickets))
	
	// Log details of closed tickets and send notification emails
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

		// TODO: Send email notification to user about ticket closure
		// This is generally a good practice to keep users informed about their tickets
		// 
		// Email details from ticket:
		// - Full Name: ticket.Name
		// - Email: ticket.Email
		// - Ticket Number: ticket.TicketNumber
		// - Subject: ticket.Subject
		//
		// Example email implementation:
		// err := j.sendTicketClosedEmail(ticket.Name, ticket.Email, ticket.TicketNumber, ticket.Subject)
		// if err != nil {
		//     log.Printf("Failed to send closure email for ticket %s: %v", ticket.TicketNumber, err)
		// }
	}

	log.Println("Auto-close support ticket job completed successfully")
}

// sendTicketClosedEmail sends an email notification when a ticket is auto-closed
// func (j *AutoCloseSupportTicket) sendTicketClosedEmail(fullName, email, ticketNumber, subject string) error {
//     // Email template for ticket closure
//     emailSubject := fmt.Sprintf("Ticket #%s has been closed", ticketNumber)
//     
//     emailBody := fmt.Sprintf(`
//         Dear %s,
//         
//         Your support ticket has been automatically closed due to inactivity.
//         
//         Ticket Details:
//         - Ticket Number: %s
//         - Subject: %s
//         - Status: Closed
//         - Closed Date: %s
//         
//         This ticket was closed because we haven't received a response from you in the past 48 hours.
//         
//         If you still need assistance with this issue, please feel free to create a new support ticket
//         or reply to this email to reopen the ticket.
//         
//         Thank you for using our service.
//         
//         Best regards,
//         Support Team
//     `, fullName, ticketNumber, subject, time.Now().Format("January 2, 2006 at 3:04 PM"))
//     
//     // Send email using your email service
//     // return j.emailService.SendEmail(email, emailSubject, emailBody)
//     
//     return nil
// }

func (j *AutoCloseSupportTicket) Schedule() string {
	return "0 0 0 * * *" // Daily at midnight
}