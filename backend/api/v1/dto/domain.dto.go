package dto

type DomainDTO struct {
	Domain string `json:"domain" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
}
