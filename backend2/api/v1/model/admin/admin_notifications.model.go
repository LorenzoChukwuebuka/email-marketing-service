package adminmodel

import "gorm.io/gorm"

type AdminNotification struct {
	gorm.Model
	UUID       string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId     string `json:"user_id"`
	Title      string `json:"title"`
	Link       string `json:"link" gorm:"default:null"`
	ReadStatus bool   `json:"read_status"`
}

type AdminNotificationResponse struct {
	ID         int     `json:"-"`
	UUID       string  `json:"uuid"`
	UserId     string  `json:"user_id"`
	Title      string  `json:"title"`
	Link       string  `json:"link"`
	ReadStatus bool    `json:"read_status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}
