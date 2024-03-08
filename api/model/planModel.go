package model

import (
	"time"
)

type Plan struct {
	ID                  int       `json:"-" gorm:"primaryKey"`
	UUID                string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	PlanName            string    `json:"planname" validate:"required" gorm:"size:255"`
	Duration            string    `json:"duration" validate:"required" gorm:"size:255"`
	Price               float32   `json:"price" validate:"required" `
	NumberOfMailsPerDay string    `json:"number_of_mails_per_day" validate:"required"`
	Details             string    `json:"details" validate:"required"`
	Status              *string   `json:"status" gorm:"size:255"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:TIMESTAMP"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt           *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type PlanResponse struct {
	ID                 int       `json:"-"`
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
