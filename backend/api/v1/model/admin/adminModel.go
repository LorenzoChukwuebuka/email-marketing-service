package adminmodel

import (
	"time"
)

type Admin struct {
	ID         int       `json:"-" gorm:"primaryKey"`
	UUID       string    `json:"uuid"`
	FirstName  *string   `json:"firstname"`
	MiddleName *string   `json:"middlename"`
	LastName   *string   `json:"lastname"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt  time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type AdminResponse struct {
	ID         int     `json:"-"`
	UUID       string  `json:"uuid"`
	FirstName  *string `json:"firstname"`
	MiddleName *string `json:"middlename"`
	LastName   *string `json:"lastname"`
	Email      string  `json:"email"`
	Password   string  `json:"-"`
	Type       string  `jsosn:"type"`
}

type AdminChangePassword struct {
	AdminId     int    `json:"admin_id" `
	OldPassword []byte `json:"old_password"`
	NewPassword []byte `json:"new_password"`
}
