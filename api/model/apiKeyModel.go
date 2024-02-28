package model

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	gorm.Model
	Id        int       `gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    int       `json:"user_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type APIKeyResponseModel struct {
	Id        int       `json:"-"`
	UUID      string    `json:"uuid"`
	UserId    int       `json:"user_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
