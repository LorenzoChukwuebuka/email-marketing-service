package model

import (
	"gorm.io/gorm"
)

type PlanStatus string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Plan struct {
	gorm.Model
	UUID                string        `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	PlanName            string        `json:"planname" gorm:"size:255"`
	Duration            string        `json:"duration" gorm:"size:255"`
	Price               float32       `json:"price"`
	NumberOfMailsPerDay string        `json:"number_of_mails_per_day"`
	Details             string        `json:"details"`
	IsPaid              bool          `json:"is_paid"`
	Status              PlanStatus    `json:"status" gorm:"size:255;default:active"`
	Features            []PlanFeature `json:"features" gorm:"foreignKey:PlanID"`
}

type PlanFeature struct {
	gorm.Model
	UUID        string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	PlanID      uint   `json:"plan_id"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	CountLimit  int    `json:"count_limit"`
	SizeLimit   int    `json:"size_limit"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
	Plan        Plan   `json:"-" gorm:"foreignKey:PlanID"`
}

type PlanResponse struct {
	ID                  uint          `json:"-"`
	UUID                string        `json:"uuid"`
	PlanName            string        `json:"planname" validate:"required"`
	Duration            string        `json:"duration" validate:"required"`
	Price               float32       `json:"price" validate:"required"`
	NumberOfMailsPerDay string        `json:"number_of_mails_per_day" validate:"required"`
	Details             string        `json:"details" validate:"required"`
	Status              PlanStatus    `json:"status"`
	Features            []PlanFeature `json:"features"`
	CreatedAt           string        `json:"created_at"`
	UpdatedAt           string        `json:"updated_at"`
	DeletedAt           *string       `json:"deleted_at,omitempty"`
}
