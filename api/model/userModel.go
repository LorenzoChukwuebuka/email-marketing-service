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
	Password   []byte       `json:"password" validate:"required"`
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
	Token    string `json:"token" validated:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePassword struct {
	UserId      int    `json:"user_id" validated:"required"`
	OldPassword []byte `json:"old_password" validated:"required"`
	NewPassword []byte `json:"new_password" validated:"required"`
}
