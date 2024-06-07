package dto

type Status string

const (
	Open     Status = "open"
	Pending  Status = "pending"
	Resolved Status = "resolved"
	Closed   Status = "closed"
)

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

type Ticket struct {
	Subject     string   `json:"subject" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Status      Status   ` json:"status" validate:"required"`
	Priority    Priority ` json:"priority" validate:"required"`
	UserID      uint     `json:"user_id"`
	Message     string   `json:"message" validate:"required"`
	TicketFile  []byte   `json:"ticket_file"`
}
