package model

import (
	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	UUID   string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId string `json:"user_id"`
	Token  string `json:"token" validated:"required"`
}
