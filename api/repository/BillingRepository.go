package repository

import (
	"email-marketing-service/api/model"
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

func (r *BillingRepository) createBillingResponse(billing *model.Billing) *model.BillingResponse {
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
		// User:          model.UserResponse{},
		// Plan:          model.PlanResponse{},
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
			UpdatedAt:           billing.Plan.UpdatedAt.Format(time.RFC3339),
			DeletedAt:           billing.Plan.DeletedAt.Format(time.RFC3339),
		}
	}

	if billing.User != nil {
		response.User = &model.UserResponse{
			UUID:       billing.User.UUID,
			FirstName:  billing.User.FirstName,
			MiddleName: billing.User.MiddleName,
			LastName:   billing.User.LastName,
		}
	}

	return &response
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

	fmt.Println(billing)

	response := r.createBillingResponse(&billing)

	return response, nil
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
