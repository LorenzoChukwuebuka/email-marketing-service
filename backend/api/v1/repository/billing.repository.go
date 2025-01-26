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

func (r *BillingRepository) CheckIfRefExists(reference string) (bool, error) {
	var count int64
	r.DB.Model(&model.Billing{}).Where("reference = ?", reference).Count(&count)
	return count > 0, nil
}

func (r *BillingRepository) createBillingResponse(billing model.Billing) model.BillingResponse {
	response := &model.BillingResponse{
		UUID:          billing.UUID,
		UserId:        billing.UserId,
		AmountPaid:    billing.AmountPaid,
		PlanId:        int(billing.PlanId),
		Duration:      billing.Duration,
		ExpiryDate:    billing.ExpiryDate.String(),
		Reference:     billing.Reference,
		TransactionId: billing.TransactionId,
		PaymentMethod: billing.PaymentMethod,
		Status:        billing.Status,
		CreatedAt:     billing.CreatedAt.String(),
		UpdatedAt:     billing.UpdatedAt.String(),
	}

	if billing.DeletedAt.Valid {
		formatted := billing.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &formatted
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
			CreatedAt:           billing.CreatedAt.String(),
			UpdatedAt:           billing.Plan.UpdatedAt.String(),
		}

		if billing.Plan.DeletedAt.Valid {
			formatted := billing.Plan.DeletedAt.Time.Format(time.RFC3339)
			response.DeletedAt = &formatted
		}

	}

	if billing.User != nil {
		response.User = &model.UserResponse{
			UUID:     billing.User.UUID,
			FullName: billing.User.FullName,
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

func mapBillingToResponse(billing model.Billing) model.BillingResponse {
	return model.BillingResponse{
		Id:            int(billing.ID),
		UUID:          billing.UUID,
		UserId:        billing.UserId,
		AmountPaid:    billing.AmountPaid,
		PlanId:        int(billing.PlanId),
		Duration:      billing.Duration,
		ExpiryDate:    billing.ExpiryDate.Format(time.RFC3339),
		Reference:     billing.Reference,
		TransactionId: billing.TransactionId,
		PaymentMethod: billing.PaymentMethod,
		Status:        billing.Status,
		CreatedAt:     billing.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     billing.UpdatedAt.Format(time.RFC3339),
		DeletedAt: func() *string {
			if billing.DeletedAt.Valid {
				str := billing.DeletedAt.Time.Format(time.RFC3339)
				return &str
			}
			return nil
		}(),
		User: func() *model.UserResponse {
			if billing.User != nil {
				return &model.UserResponse{
					ID:          billing.User.ID,
					UUID:        billing.User.UUID,
					FullName:    billing.User.FullName,
					Email:       billing.User.Email,
					Company:     billing.User.Company,
					PhoneNumber: billing.User.PhoneNumber,
					Password: func() *string {
						if billing.User != nil && billing.User.Password != nil {
							buser := billing.User.Password
							return buser
						}
						return nil
					}(),
					Verified:  billing.User.Verified,
					Blocked:   billing.User.Blocked,
					CreatedAt: billing.User.CreatedAt.Format(time.RFC3339),

					UpdatedAt: billing.User.UpdatedAt.Format(time.RFC3339),
					DeletedAt: func() *string {
						if billing.User.DeletedAt.Valid {
							str := billing.User.DeletedAt.Time.Format(time.RFC3339)
							return &str
						}
						return nil
					}(),
				}
			}
			return nil
		}(),
		Plan: func() *model.PlanResponse {
			if billing.Plan != nil {
				return &model.PlanResponse{
					ID:                  billing.Plan.ID,
					UUID:                billing.Plan.UUID,
					PlanName:            billing.Plan.PlanName,
					Duration:            billing.Plan.Duration,
					Price:               billing.Plan.Price,
					NumberOfMailsPerDay: billing.Plan.NumberOfMailsPerDay,
					Details:             billing.Plan.Details,
					Status:              billing.Plan.Status,
					Features:            billing.Plan.Features,
					CreatedAt:           billing.Plan.CreatedAt.Format(time.RFC3339),
					UpdatedAt:           billing.Plan.UpdatedAt.Format(time.RFC3339),
					DeletedAt: func() *string {
						if billing.Plan.DeletedAt.Valid {
							str := billing.Plan.DeletedAt.Time.Format(time.RFC3339)
							return &str
						}
						return nil
					}(),
				}
			}
			return nil
		}(),
	}
}

func (r *BillingRepository) GetAllPayments(userID int, params PaginationParams) (PaginatedResult, error) {
	var billingRecords []model.Billing
	query := r.DB.
		Preload("User").
		Preload("Plan").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&billingRecords)

	paginatedResult, err := Paginate(query, params, &billingRecords)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	var billingResponses []model.BillingResponse
	for _, billing := range billingRecords {
		billingResponses = append(billingResponses, mapBillingToResponse(billing))
	}

	paginatedResult.Data = billingResponses

	return paginatedResult, nil
}
