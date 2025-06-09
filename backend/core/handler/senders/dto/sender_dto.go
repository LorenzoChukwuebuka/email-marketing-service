package dto

import "time"

type SenderDTO struct {
	UserID    string `json:"user_id" validate:"required"`
	CompanyID string `json:"company_id" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Name      string `json:"name" validate:"required"`
	SenderId  string `json:"sender_id"`
}

type VerifySenderDTO struct {
	UserID    string `json:"user_id" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Token     string `json:"token" validate:"required"`
	CompanyID string `json:"company_id"`
}

type FetchSenderDTO struct {
	UserID      string `json:"user_id"  `
	CompanyID   string `json:"company_id"  `
	SenderId    string `json:"sender_id"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	SearchQuery string `json:"search_query"`
}

type SendersResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	CompanyID string     `json:"company_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Verified  bool       `json:"verified"`
	IsSigned  bool       `json:"is_signed"`
	DomainID  string     `json:"domain_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
