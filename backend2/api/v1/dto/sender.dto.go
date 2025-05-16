package dto

type SenderDTO struct {
	UserID   string `json:"user_id" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	SenderId string `json:"sender_id"`
}

type VerifySenderDTO struct {
	UserID string `json:"user_id" validate:"required"`
	Email  string `json:"email" validate:"required"`
	Token  string `json:"token" validate:"required"`
}
