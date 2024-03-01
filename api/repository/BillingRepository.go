package repository

import (
	"email-marketing-service/api/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BillingRepository struct {
	DB *gorm.DB
}

func NewBillingRepository(db *gorm.DB) *BillingRepository {
	return &BillingRepository{DB: db}
}

func (r *BillingRepository) createBillingResponse(billing model.Billing) model.BillingResponse {
	response := model.BillingResponse{
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
		UpdatedAt:     billing.UpdatedAt.Format(time.RFC3339),
		DeletedAt:     billing.DeletedAt.Format(time.RFC3339),
		User:          model.UserResponse{},
		Plan:          model.PlanResponse{},
	}

	if billing.Plan.ID != 0 {
		response.Plan = model.PlanResponse{
			UUID:                billing.Plan.UUID,
			PlanName:            billing.Plan.PlanName,
			Duration:            billing.Plan.Duration,
			Price:               billing.Plan.Price,
			NumberOfMailsPerDay: billing.Plan.NumberOfMailsPerDay,
			Details:             billing.Plan.Details,
			Status:              billing.Plan.Status,
			CreatedAt:           billing.Plan.CreatedAt,
			UpdatedAt:           billing.Plan.UpdatedAt.Format(time.RFC3339),
			DeletedAt:           billing.Plan.DeletedAt.Format(time.RFC3339),
		}
	}

	if billing.User.ID != 0 {
		response.User = model.UserResponse{
			UUID:       billing.User.UUID,
			FirstName:  billing.User.FirstName,
			MiddleName: billing.User.MiddleName,
			LastName:   billing.User.LastName,
		}
	}

	return response
}

func (r *BillingRepository) CreateBilling(d *model.Billing) (*model.Billing, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert plan: %w", err)
	}
	return d, nil
}

func (r *BillingRepository) GetSingleBillingRecord(billingID string, userID int) (*model.BillingResponse, error) {

	var billingResponse model.Billing

	// Query the database and preload associated data
	if err := r.DB.
		Preload("User").
		Preload("Plan").
		Where("uuid = ? AND user_id = ?", billingID, userID).
		First(&billingResponse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("billing record not found")
		}
		return nil, fmt.Errorf("failed to get billing record: %w", err)
	}

	response := r.createBillingResponse(billingResponse)

	return &response, nil
}

func (r *BillingRepository) GetAllPayments(userID int, page int) ([]model.BillingResponse, error) {
	pageSize := 20

	offset := (page - 1) * pageSize

	var billingRecords []model.BillingResponse
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
