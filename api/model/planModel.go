package model

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	Id                  int       `json:"-" gorm:"primaryKey"`
	UUID                string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	PlanName            string    `json:"planname" validate:"required"`
	Duration            string    `json:"duration" validate:"required"`
	Price               float32   `json:"price" validate:"required"`
	NumberOfMailsPerDay string    `json:"number_of_mails_per_day" validate:"required"`
	Details             string    `json:"details" validate:"required"`
	Status              *string   `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt           time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type PlanResponse struct {
	Id                  int       `json:"-"`
	UUID                string    `json:"uuid"`
	PlanName            string    `json:"planname" validate:"required"`
	Duration            string    `json:"duration" validate:"required"`
	Price               float32   `json:"price" validate:"required"`
	NumberOfMailsPerDay string    `json:"number_of_mails_per_day" validate:"required"`
	Details             string    `json:"details" validate:"required"`
	Status              *string   `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           string    `json:"updated_at"`
	DeletedAt           string    `json:"deleted_at"`
}
