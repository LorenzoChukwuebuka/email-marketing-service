package model

import (
	"time"
)

type APIKey struct {
	Id        int       `json:"-" gorm:"primaryKey;index"`
	UUID      string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
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
