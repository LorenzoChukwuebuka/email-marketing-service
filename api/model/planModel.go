package model

import (
	"database/sql"
	"time"
)

type PlanModel struct {
	Id        int          `json:"id"`
	PlanName  string       `json:"planname" validate:"required"`
	Duration  string       `json:"duration" validate:"required"`
	Price     float32      `json:"price" validate:"required"`
	Details   string       `json:"lastname" validate:"required"`
	Status    *string      `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
