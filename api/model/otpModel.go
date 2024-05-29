package model

import "time"

type OTP struct {
	Id        int       `gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    string       `json:"user_id"`
	Token     string    `json:"token" validated:"required"`
	CreatedAt time.Time `json:"created_at"`
}

