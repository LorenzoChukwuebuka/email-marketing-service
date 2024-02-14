package model

import (
	"time"
)

type User struct {
	ID         int       `gorm:"primaryKey"`
	UUID       string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	FirstName  string    `json:"firstname" validate:"required"`
	MiddleName *string   `json:"middlename"`
	LastName   string    `json:"lastname" validate:"required"`
	UserName   string    `json:"username" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	VerifiedAt time.Time `json:"verified_at" gorm:"type:TIMESTAMP;null;default:null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt  time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type LoginModel struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
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
	OldPassword string `json:"old_password" validated:"required"`
	NewPassword string `json:"new_password" validated:"required"`
}

type UserResponse struct {
	ID         int       `json:"-"`
	UUID       string    `json:"uuid"`
	FirstName  string    `json:"firstname"`
	MiddleName *string   `json:"middlename"`
	LastName   string    `json:"lastname"`
	UserName   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	VerifiedAt string    `json:"verified_at"`
	UpdatedAt  string    `json:"updated_at"`
	DeletedAt  string    `json:"deleted_at"`
}
