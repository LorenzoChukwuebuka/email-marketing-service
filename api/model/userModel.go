package model

import "time"

type User struct {
	FirstName  string    `json:"firstname" validate:"required"`
	MiddleName string    `json:"middlename"`
	LastName   string    `json:"lastname" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   []byte    `json:"password" validate:"required"`
	VerifiedAt time.Time `json:"verified_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
