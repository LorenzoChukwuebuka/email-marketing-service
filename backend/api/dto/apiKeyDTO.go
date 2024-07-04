package dto

type APIkeyDTO struct {
	UserId string `json:"user_id"`
	Name string `json:"name" validate:"required"`
}
