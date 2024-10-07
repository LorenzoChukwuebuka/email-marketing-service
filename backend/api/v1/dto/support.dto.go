package dto

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
)

type SupportStatus string

const (
	OpenTicket     SupportStatus = "open"
	CloseTicket    SupportStatus = "closed"
	ResolvedTicket SupportStatus = "resolved"
	PendingTicket  SupportStatus = "pending"
)

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

type CreateSupportTicketRequest struct {
	Subject     string                  `form:"subject" validate:"required,max=255"`
	Description string                  `form:"description" `
	Priority    Priority                `form:"priority" validate:"required,validpriority"`
	Message     string                  `form:"message" validate:"required"`
	File        []*multipart.FileHeader `form:"file"`
}

type CreateSupportTicketResponse struct {
	TicketID uint   `json:"ticket_id"`
	Message  string `json:"message"`
}

type ReplyTicketRequest struct {
	Message string                  `form:"message" validate:"required"`
	File    []*multipart.FileHeader `form:"file"`
}

type ReplyTicketResponse struct {
	MessageID uint   `json:"message_id"`
	Message   string `json:"message"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("validpriority", isValidPriority)
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

func ValidateCreateSupportTicketRequest(req *CreateSupportTicketRequest) error {
	return validate.Struct(req)
}

func ValidateReplyTicketRequest(req *ReplyTicketRequest) error {
	return validate.Struct(req)
}
