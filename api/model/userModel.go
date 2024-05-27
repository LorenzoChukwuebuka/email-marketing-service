package model

import (
	"time"
)

type User struct {
	ID         int       `json:"-" gorm:"primaryKey"`
	UUID       string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FullName   string    `json:"fullname" validate:"required"`
	Company    string    `json:"company" validate:"required"`
	Email      string    `json:"email" validate:"required,email" gorm:"index"`
	Password   string    `json:"password" validate:"required" gorm:"index"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	VerifiedAt time.Time `json:"verified_at" gorm:"type:TIMESTAMP;null;default:null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt  time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}



type UserResponse struct {
	ID         int       `json:"-"`
	UUID       string    `json:"uuid,omitempty"`
	FullName   string    `json:"fullname,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	Verified   bool      `json:"verified,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	VerifiedAt string    `json:"verified_at,omitempty"`
	UpdatedAt  string    `json:"updated_at,omitempty"`
	DeletedAt  string    `json:"deleted_at,omitempty"`
}
