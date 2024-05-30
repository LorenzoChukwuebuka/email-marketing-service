package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type DailyMailCalcRepository struct {
	DB *gorm.DB
}

func NewDailyMailCalcRepository(db *gorm.DB) *DailyMailCalcRepository {
	return &DailyMailCalcRepository{DB: db}
}

func (r *DailyMailCalcRepository) CreateRecordDailyMailCalculation(d *model.DailyMailCalc) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *DailyMailCalcRepository) GetDailyMailRecordForToday(subscriptionId int) (*model.DailyMailCalcResponseModel, error) {
	var record model.DailyMailCalc

	today := time.Now().Format("2006-01-02")
	startOfDay, _ := time.Parse("2006-01-02 15:04:05", today+" 00:00:00")
	endOfDay, _ := time.Parse("2006-01-02 15:04:05", today+" 23:59:59")

	if err := r.DB.Where("subscription_id = ? AND created_at >= ? AND created_at <= ?", subscriptionId, startOfDay, endOfDay).First(&record).Error; err != nil {
		return nil, fmt.Errorf("error fetching record: %w", err)
	}

	response := &model.DailyMailCalcResponseModel{
		ID:             record.ID,
		UUID:           record.UUID,
		SubscriptionID: record.SubscriptionID,
		MailsForADay:   record.MailsForADay,
		MailsSent:      record.MailsSent,
		RemainingMails: record.RemainingMails,
		CreatedAt:      record.CreatedAt,
	}

	if record.UpdatedAt != nil {
		updatedAt := record.UpdatedAt.Format(time.RFC3339)
		response.UpdatedAt = &updatedAt
	}

	fmt.Printf("Record: %+v\n", record)

	return response, nil
}

func (r *DailyMailCalcRepository) UpdateDailyMailCalcRepository(d *model.DailyMailCalc) error {
	// Fetch the existing record from the database
	var existingRecord model.DailyMailCalc
	if err := r.DB.Where("uuid = ?", d.UUID).First(&existingRecord).Error; err != nil {
		return fmt.Errorf("error fetching daily mail record: %w", err)
	}

	// Update the fields
	existingRecord.RemainingMails = d.RemainingMails
	existingRecord.MailsSent = d.MailsSent

	// Save the updated record back to the database
	if err := r.DB.Save(&existingRecord).Error; err != nil {
		return fmt.Errorf("error updating daily mail record: %w", err)
	}

	return nil
}
