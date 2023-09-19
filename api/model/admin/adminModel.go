package adminmodel

import (
	"database/sql"
	"time"
)

type AdminModel struct {
	ID         int          `json:"id"`
	UUID        string       `json:"uuid"`
	FirstName  *string      `json:"firstname"`
	MiddleName *string      `json:"middlename"`
	LastName   *string      `json:"lastname"`
	Email      string       `json:"email"`
	Password   []byte       `json:"password"`
	Type       string       `json:"type"`
	CreatedAt  time.Time    `json:"created_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}

type AdminLogin struct {
	Email    string `json:"email" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type AdminResponse struct {
	ID         int     `json:"id"`
	UUID        string       `json:"uuid"`
	FirstName  *string `json:"firstname"`
	MiddleName *string `json:"middlename"`
	LastName   *string `json:"lastname"`
	Email      string  `json:"email"`
	Password   []byte  `json:"password"`
	Type       string  `jsosn:"type"`
}

type AdminChangePassword struct {
	AdminId     int    `json:"admin_id" `
	OldPassword []byte `json:"old_password"`
	NewPassword []byte `json:"new_password"`
}
