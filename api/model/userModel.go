package model

import "time"

type User struct {
	UUID       string    `json:"uuid"`
	FirstName  string    `json:"firstname" validate:"required"`
	MiddleName string    `json:"middlename"`
	LastName   string    `json:"lastname" validate:"required"`
	UserName   string    `json:"username" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   []byte    `json:"password" validate:"required"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	VerifiedAt time.Time `json:"verified_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
