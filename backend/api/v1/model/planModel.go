package model

import (
	"gorm.io/gorm"
)

type status string

const (
	Active   status = "active"
	Inactive status = "inactive"
)

type Plan struct {
	gorm.Model
	UUID                string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	PlanName            string     `json:"planname" gorm:"size:255"`
	Duration            string     `json:"duration" gorm:"size:255"`
	Price               float32    `json:"price" `
	NumberOfMailsPerDay string     `json:"number_of_mails_per_day"`
	Details             string     `json:"details"`
	Status              *string    `json:"status" gorm:"size:255;default:active"`
	
}

type PlanResponse struct {
	ID                  uint       `json:"-"`
	UUID                string    `json:"uuid"`
	PlanName            string    `json:"planname" validate:"required"`
	Duration            string    `json:"duration" validate:"required"`
	Price               float32   `json:"price" validate:"required"`
	NumberOfMailsPerDay string    `json:"number_of_mails_per_day" validate:"required"`
	Details             string    `json:"details" validate:"required"`
	Status              *string   `json:"status"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string    `json:"updated_at"`
	DeletedAt           *string    `json:"deleted_at"`
}
