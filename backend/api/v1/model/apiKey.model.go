package model

import (
	"gorm.io/gorm"
)

type KeyStatus string

const (
	KeyActive   KeyStatus = "active"
	KeyInactive Status    = "inactive"
)

type APIKey struct {
	gorm.Model
	UUID   string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId string `json:"user_id" gorm:"type:uuid"`
	Name   string `json:"name"`
	APIKey string `json:"api_key" gorm:"index"`
}

type APIKeyResponseModel struct {
	Id        int    `json:"-"`
	UUID      string `json:"uuid,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	Name      string `json:"name,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type SMTPMasterKey struct {
	gorm.Model
	UUID      string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId    string    `json:"user_id" gorm:"type:uuid"`
	SMTPLogin string    `json:"smtp_login"`
	KeyName   string    `json:"key_name"`
	Password  string    `json:"password"`
	Status    KeyStatus `json:"status"`
}

type SMTPKey struct {
	gorm.Model
	UUID     string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId   string    `json:"user_id" gorm:"type:uuid"`
	KeyName  string    `json:"key_name" gorm:"index"`
	Password string    `json:"password"`
	Status   KeyStatus `json:"status"`
}

type SMTPDetailsResponse struct {
	Id        int     `json:"-" `
	UUID      string  `json:"uuid"`
	UserId    string  `json:"user_id"`
	KeyName   string  `json:"key_name" `
	SMTPLogin string  `json:"smtp_login"`
	Password  string  `json:"password"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at" `
	DeletedAt *string `json:"deleted_at"`
}
