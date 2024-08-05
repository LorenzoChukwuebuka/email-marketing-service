package dto

import (
	"github.com/go-playground/validator/v10"
)

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

type SupportTicket struct {
	Subject     string   `json:"subject" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Status      Status   `json:"status" validate:"required,validstatus"`
	Priority    Priority `json:"priority" validate:"required,validpriority"`
	UserID      string   `json:"user_id"`
	Message     string   `json:"message" validate:"required"`
	TicketFile  *[]byte   `json:"ticket_file"`
	FilePath    string
}

func isValidStatus(fl validator.FieldLevel) bool {
	status := Status(fl.Field().String())
	switch status {
	case Open, Pending, Resolved, Closed:
		return true
	default:
		return false
	}
}

func isValidPriority(fl validator.FieldLevel) bool {
	priority := Priority(fl.Field().String())
	switch priority {
	case Low, Medium, High:
		return true
	default:
		return false
	}
}

var validate *validator.Validate

func InitValidate() {
	validate = validator.New()
	validate.RegisterValidation("validstatus", isValidStatus)
	validate.RegisterValidation("validpriority", isValidPriority)
}
