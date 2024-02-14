package model

import (
	"time"
)

type UserSession struct {
	Id        int       `gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    int       `json:"user_id"`
	Device    *string   `json:"device"`
	IPAddress *string   `json:"ip_address"`
	Browser   *string   `json:"browser"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type UserSessionResponseModel struct {
	Id        int       `json:"-"`
	UUID      string    `json:"uuid"`
	UserId    int       `json:"user_id"`
	Device    *string   `json:"device"`
	IPAddress *string   `json:"ip_address"`
	Browser   *string   `json:"browser"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	DeletedAt string    `json:"deleted_at"`
}
