package model

import (
	"database/sql"
	"time"
)

type PlanModel struct {
	Id                  int          `json:"id"`
	UUID                string       `json:"uuid"`
	PlanName            string       `json:"planname" validate:"required"`
	Duration            string       `json:"duration" validate:"required"`
	Price               float32      `json:"price" validate:"required"`
	NumberOfMailsPerDay string       `json:"number_of_mails_per_day" validate:"required"`
	Details             string       `json:"details" validate:"required"`
	Status              *string      `json:"status"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           sql.NullTime `json:"updated_at"`
	DeletedAt           sql.NullTime `json:"deleted_at"`
}

type PlanResponse struct {
	Id                  int       `json:"id"`
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
