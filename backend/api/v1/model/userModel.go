package model

import (
	"time"
)

type User struct {
	ID          int       `json:"-" gorm:"primaryKey"`
	UUID        string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FullName    string    `json:"fullname" `
	Company     string    `json:"company" `
	Email       string    `json:"email"  gorm:"index"`
	PhoneNumber string    `json:"phonenumber" gorm:"type:varchar(255);default:null"`
	Password    string    `json:"password"  gorm:"index"`
	Verified    bool      `json:"verified"`
	Blocked     bool      `json:"blocked" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	VerifiedAt  time.Time `json:"verified_at" gorm:"type:TIMESTAMP;null;default:null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type UserResponse struct {
	ID          int       `json:"-"`
	UUID        string    `json:"uuid,omitempty"`
	FullName    string    `json:"fullname,omitempty"`
	Email       string    `json:"email,omitempty"`
	Company     string    `json:"company"`
	PhoneNumber string    `json:"phonenumber,omitempty"`
	Password    string    `json:"-"`
	Verified    bool      `json:"verified,omitempty"`
	Blocked     bool      `json:"blocked"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	VerifiedAt  string    `json:"verified_at,omitempty"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
	DeletedAt   string    `json:"deleted_at,omitempty"`
}
