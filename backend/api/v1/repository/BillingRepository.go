package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
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

func (r *BillingRepository) CreateBilling(d *model.Billing) (*model.Billing, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert plan: %w", err)
	}
	return d, nil
}

func (r *BillingRepository) createBillingResponse(billing model.Billing) model.BillingResponse {
	response := &model.BillingResponse{
		UUID:          billing.UUID,
		UserId:        billing.UserId,
		AmountPaid:    billing.AmountPaid,
		PlanId:        billing.PlanId,
		Duration:      billing.Duration,
		ExpiryDate:    billing.ExpiryDate,
		Reference:     billing.Reference,
		TransactionId: billing.TransactionId,
		PaymentMethod: billing.PaymentMethod,
		Status:        billing.Status,
		CreatedAt:     billing.CreatedAt,
	}

	if billing.UpdatedAt != nil {
		response.UpdatedAt = billing.UpdatedAt.Format(time.RFC3339)
	} else {
		response.UpdatedAt = ""
	}

	if billing.DeletedAt != nil {
		response.DeletedAt = billing.DeletedAt.Format(time.RFC3339)
	} else {
		response.UpdatedAt = ""
	}

	if billing.Plan != nil {
		response.Plan = &model.PlanResponse{
			UUID:                billing.Plan.UUID,
			PlanName:            billing.Plan.PlanName,
			Duration:            billing.Plan.Duration,
			Price:               billing.Plan.Price,
			NumberOfMailsPerDay: billing.Plan.NumberOfMailsPerDay,
			Details:             billing.Plan.Details,
			Status:              billing.Plan.Status,
			CreatedAt:           billing.Plan.CreatedAt,
		}

		if billing.Plan.UpdatedAt != nil {
			response.Plan.UpdatedAt = billing.Plan.UpdatedAt.Format(time.RFC3339)
		} else {
			response.Plan.UpdatedAt = ""
		}

		if billing.DeletedAt != nil {
			response.DeletedAt = billing.DeletedAt.Format(time.RFC3339)
		} else {
			response.UpdatedAt = ""
		}

	}

	if billing.User != nil {
		response.User = &model.UserResponse{
			UUID:       billing.User.UUID,
			FullName:  billing.User.FullName,
			 
		}
	}

	return *response
}

func (r *BillingRepository) GetSingleBillingRecord(billingID string, userID int) (*model.BillingResponse, error) {

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

	response := r.createBillingResponse(billing)

	return &response, nil
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
