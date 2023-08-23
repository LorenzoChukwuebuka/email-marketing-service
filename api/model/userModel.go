package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int            `json:"id"`
	UUID       string         `json:"uuid"`
	FirstName  string         `json:"firstname" validate:"required"`
	MiddleName sql.NullString `json:"middlename"`
	LastName   string         `json:"lastname" validate:"required"`
	UserName   string         `json:"username" validate:"required"`
	Email      string         `json:"email" validate:"required,email"`
	Password   []byte         `json:"password" validate:"required"`
	Verified   bool           `json:"verified"`
	CreatedAt  time.Time      `json:"created_at"`
	VerifiedAt sql.NullTime   `json:"verified_at"`
	UpdatedAt  sql.NullTime   `json:"updated_at"`
	DeletedAt  sql.NullTime   `json:"deleted_at"`
}
