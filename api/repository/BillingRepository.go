package repository

import (
	"database/sql"
)

type BillingRepository struct {
	DB *sql.DB
}

func NewBillingRepository(db *sql.DB) *BillingRepository {
	return &BillingRepository{DB: db}
}
