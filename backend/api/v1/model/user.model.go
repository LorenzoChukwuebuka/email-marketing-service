package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UUID        string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FullName    string     `json:"fullname" `
	Company     string     `json:"company" `
	Email       string     `json:"email"  gorm:"index"`
	PhoneNumber string     `json:"phonenumber" gorm:"type:varchar(255);default:null"`
	Password    string     `json:"password"  gorm:"index"`
	Verified    bool       `json:"verified"`
	Blocked     bool       `json:"blocked" gorm:"default:false"`
	VerifiedAt  *time.Time `json:"verified_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type UserResponse struct {
	ID          uint    `json:"-"`
	UUID        string  `json:"uuid"`
	FullName    string  `json:"fullname"`
	Email       string  `json:"email"`
	Company     string  `json:"company"`
	PhoneNumber string  `json:"phonenumber"`
	Password    string  `json:"-"`
	Verified    bool    `json:"verified"`
	Blocked     bool    `json:"blocked"`
	CreatedAt   string  `json:"created_at"`
	VerifiedAt  *string `json:"verified_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at"`
}

type UserTempEmail struct {
	gorm.Model
	TemporaryEmail string `gorm:"unique;not null"`
	UserId         string
}
