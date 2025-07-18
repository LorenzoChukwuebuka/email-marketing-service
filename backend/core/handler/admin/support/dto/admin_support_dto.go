package dto

import (
	"time"
)

type SupportTicketResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Subject      string     `json:"subject"`
	Description  *string    `json:"description"`
	TicketNumber string     `json:"ticket_number"`
	Status       *string    `json:"status"`
	Priority     *string    `json:"priority"`
	LastReply    *time.Time `json:"last_reply"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
