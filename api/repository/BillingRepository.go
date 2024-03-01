package repository

import (
	"email-marketing-service/api/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type BillingRepository struct {
	DB *gorm.DB
}

func NewBillingRepository(db *gorm.DB) *BillingRepository {
	return &BillingRepository{DB: db}
}

func (r *BillingRepository) CreateBilling(d *model.Billing) (*model.Billing, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert plan: %w", err)
	}
	return d, nil
}

func (r *BillingRepository) GetSingleBillingRecord(billingID string, userID int) (*model.Billing, error) {

	var billing model.Billing

	// Query the database and preload associated data
	if err := r.DB.
		Preload("User").
		Preload("Plan").
		Where("uuid = ? AND user_id = ?", billingID, userID).
		First(&billing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("billing record not found")
		}
		return nil, fmt.Errorf("failed to get billing record: %w", err)
	}

	fmt.Println(billing)

	return &billing, nil
}

func (r *BillingRepository) GetAllPayments(userID int, page int) ([]model.Billing, error) {
	pageSize := 20

	offset := (page - 1) * pageSize

	var billingRecords []model.Billing
	if err := r.DB.
		Preload("User").
		Preload("Plan").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&billingRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to get billing records: %w", err)
	}

	return billingRecords, nil
}
