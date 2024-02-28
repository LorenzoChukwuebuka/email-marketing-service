package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BillingRepository struct {
	DB *gorm.DB
}

func NewBillingRepository(db *gorm.DB) *BillingRepository {
	return &BillingRepository{DB: db}
}

var SetTime = func(field sql.NullTime, target *string) {
	if field.Valid {
		*target = field.Time.Format(time.RFC3339Nano)
	}
}

func (r *BillingRepository) CreateBilling(d *model.Billing) (*model.Billing, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert plan: %w", err)
	}
	return d, nil
}

func (r *BillingRepository) GetSingleBillingRecord(billingID string, userID int) (*model.BillingResponse, error) {
	return nil, nil
}

func (r *BillingRepository) GetAllPayments(userId int, page int) ([]model.BillingResponse, error) {

	// Assuming a fixed page size of 20
	pageSize := 20

	// Calculate the offset based on the page number and fixed page size
	offset := (page - 1) * pageSize

	print(offset)

	return nil, nil
}
