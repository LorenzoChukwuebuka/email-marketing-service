package adminrepository

import (
	"email-marketing-service/api/v1/model"
	"gorm.io/gorm"
	"time"
)

type AdminBillingRepository struct {
	DB *gorm.DB
}

func NewAdminBillingRepository(db *gorm.DB) *AdminBillingRepository {
	return &AdminBillingRepository{
		DB: db,
	}
}

func (r *AdminBillingRepository) GetTotalBillingForUser(userId uint) (float32, error) {
	var totalBilling float32
	err := r.DB.Model(&model.Billing{}).
		Where("user_id = ?", userId).
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&totalBilling).Error
	if err != nil {
		return 0, err
	}
	return totalBilling, nil

}

func (r *AdminBillingRepository) GetAllBillingsForUser(userId uint) ([]model.Billing, error) {
	var billings []model.Billing
	err := r.DB.Preload("User").
		Where("user_id = ?", userId).
		Find(&billings).Error
	if err != nil {
		return nil, err
	}
	return billings, nil
}

func (r *AdminBillingRepository) GetTotalBillingsByTimeRange(startDate, endDate time.Time) (float32, error) {
	var totalBilling float32
	err := r.DB.Model(&model.Billing{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&totalBilling).Error
	if err != nil {
		return 0, err
	}
	return totalBilling, nil
}
