package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MonthlyMailCalcRepository struct {
	DB *gorm.DB
}

func NewMonthlyMailCalcRepository(db *gorm.DB) *MonthlyMailCalcRepository {
	return &MonthlyMailCalcRepository{DB: db}
}

func (r *MonthlyMailCalcRepository) CreateRecordMonthlyMailCalculation(d *model.MonthlyMailCalc) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *MonthlyMailCalcRepository) UpdateMonthlyMailCalcRepository(d *model.MonthlyMailCalc) error {
	// Fetch the existing record from the database
	var existingRecord model.MonthlyMailCalc
	if err := r.DB.Where("uuid = ?", d.UUID).First(&existingRecord).Error; err != nil {
		return fmt.Errorf("error fetching Monthly mail record: %w", err)
	}

	// Update the fields
	existingRecord.RemainingMails = d.RemainingMails
	existingRecord.MailsSent = d.MailsSent

	// Save the updated record back to the database
	if err := r.DB.Save(&existingRecord).Error; err != nil {
		return fmt.Errorf("error updating Monthly mail record: %w", err)
	}

	return nil
}

func (r *MonthlyMailCalcRepository) GetUserActiveCalculation(subscriptionId int) (*model.MonthlyMailCalcResponseModel, error) {
	var record model.MonthlyMailCalc

	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour).Add(-time.Nanosecond)

	err := r.DB.Where("subscription_id = ? AND created_at BETWEEN ? AND ?", subscriptionId, todayStart, todayEnd).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	response := &model.MonthlyMailCalcResponseModel{
		ID:             int(record.ID),
		UUID:           record.UUID,
		SubscriptionID: record.SubscriptionID,
		MailsForAMonth: record.MailsForAMonth,
		MailsSent:      record.MailsSent,
		RemainingMails: record.RemainingMails,
		CreatedAt:      record.CreatedAt.String(),
	}

	return response, nil
}
