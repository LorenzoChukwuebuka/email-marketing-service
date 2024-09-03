package model

import "gorm.io/gorm"

type Sender struct {
	gorm.Model
	UUID   string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Domains struct {
	gorm.Model
	UserID string `json:"user_id"`
	Domain string `json:"domain"`
}
