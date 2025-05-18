package dto

type APIkeyRequestDTO struct {
	UserId string `json:"user_id"`
	CompanyID string `json:"company_id"`
	Name   string `json:"name" validate:"required"`
}

type SMTPKeyRequestDTO struct {
	UserId  string `json:"user_id"`
	CompanyID string `json:"company_id"`
	KeyName string `json:"key_name" validate:"required"`
}
