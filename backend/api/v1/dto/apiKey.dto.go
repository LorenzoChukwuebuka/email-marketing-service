package dto

type APIkeyDTO struct {
	UserId string `json:"user_id"`
	Name   string `json:"name" validate:"required"`
}

type SMTPKeyDTO struct {
	UserId  string `json:"user_id"`
	KeyName string `json:"key_name" validate:"required"`
}
