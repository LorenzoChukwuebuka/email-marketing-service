package model

import (
	"time"
)

type APIKey struct {
	Id        int       `json:"-" gorm:"primaryKey;index"`
	UUID      string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId    string    `json:"user_id"`
	APIKey    string    `json:"api_key" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type APIKeyResponseModel struct {
	Id        int       `json:"-"`
	UUID      string    `json:"uuid"`
	UserId    string    `json:"user_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
