package dto

import (
	"email-marketing-service/internal/enums"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("validpriority", isValidPriority)
}


type CreateSupportTicketRequest struct {
	Subject     string                  `form:"subject" validate:"required,max=255"`
	Description string                  `form:"description" `
	Priority    enums.Priority          `form:"priority" validate:"required,validpriority"`
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


func isValidPriority(fl validator.FieldLevel) bool {
	priority := enums.Priority(fl.Field().String())
	switch priority {
	case enums.Low, enums.Medium, enums.High:
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

type SupportTicketResponse struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Subject       string          `json:"subject"`
	Description   *string         `json:"description"`
	TicketNumber  string          `json:"ticket_number"`
	Status        *string         `json:"status"`
	Priority      *string         `json:"priority"`
	LastReply     *time.Time      `json:"last_reply"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	TicketFile    []TicketFile    `json:"ticket_files"`
	TicketMessage []TicketMessage `json:"ticket_messages"`
}

type TicketFile struct {
	ID        string `json:"id"`
	MessageID string `json:"message_id"`
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
}

type TicketMessage struct {
	ID        string         `json:"id"`
	TicketID  string         `json:"ticket_id"`
	UserID    string         `json:"user_id"`
	Message   string         `json:"message"`
	IsAdmin   bool           `json:"is_admin"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	User      *UserResponse  `json:"user,omitempty"`
	Admin     *AdminResponse `json:"admin,omitempty"`
}

type UserResponse struct {
	ID         string     `json:"id"`
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	Verified   bool       `json:"verified"`
	Blocked    bool       `json:"blocked"`
	VerifiedAt *time.Time `json:"verified_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type AdminResponse struct {
	ID         string     `json:"id"`
	Firstname  string     `json:"firstname"`
	Middlename string     `json:"middlename,omitempty"`
	Lastname   string     `json:"lastname"`
	Email      string     `json:"email"`
	Type       string     `json:"type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
