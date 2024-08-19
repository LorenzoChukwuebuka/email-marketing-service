package model

import (
	"gorm.io/gorm"
	 
)

type UserSession struct {
	gorm.Model
	UUID      string  `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    string  `json:"user_id"`
	Device    *string `json:"device"`
	IPAddress *string `json:"ip_address"`
	Browser   *string `json:"browser"`
}

type UserSessionResponseModel struct {
	Id        int       `json:"-"`
	UUID      string    `json:"uuid"`
	UserId    string    `json:"user_id"`
	Device    *string   `json:"device"`
	IPAddress *string   `json:"ip_address"`
	Browser   *string   `json:"browser"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	DeletedAt string    `json:"deleted_at"`
}
