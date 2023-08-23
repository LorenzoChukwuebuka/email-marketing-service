package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int          `json:"id"`
	UUID       string       `json:"uuid"`
	FirstName  string       `json:"firstname" validate:"required"`
	MiddleName *string      `json:"middlename"`
	LastName   string       `json:"lastname" validate:"required"`
	UserName   string       `json:"username" validate:"required"`
	Email      string       `json:"email" validate:"required,email"`
	Password   []byte       `json:"password" validate:"required,min=8,max=50"`
	Verified   bool         `json:"verified"`
	CreatedAt  time.Time    `json:"created_at"`
	VerifiedAt sql.NullTime `json:"verified_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}

type LoginModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password []byte `json:"password" validate:"required"`
}

type ForgetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Password []byte `json:"password" validate:"required,min=8,max=50"`
	Token    string `json:"token" validate:"required"`
}
