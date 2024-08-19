package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MailUsageRepository struct {
	DB *gorm.DB
}

func NewMailUsageRepository(db *gorm.DB) *MailUsageRepository {
	return &MailUsageRepository{DB: db}
}

func (r *MailUsageRepository) CreateMailUsageRecord(m *model.MailUsage) error {
	if err := r.DB.Create(m).Error; err != nil {
		return fmt.Errorf("failed to insert mail usage record: %w", err)
	}
	return nil
}

func (r *MailUsageRepository) GetCurrentMailUsageRecord(subscriptionId int) (*model.MailUsageResponseModel, error) {
	var record model.MailUsage

	now := time.Now()

	if err := r.DB.Where("subscription_id = ? AND period_start <= ? AND period_end >= ?", subscriptionId, now, now).First(&record).Error; err != nil {
		return nil, fmt.Errorf("error fetching record: %w", err)
	}

	response := &model.MailUsageResponseModel{
		ID:             int(record.ID),
		UUID:           record.UUID,
		SubscriptionID: record.SubscriptionID,
		PeriodStart:    record.PeriodStart,
		PeriodEnd:      record.PeriodEnd,
		LimitAmount:    record.LimitAmount,
		MailsSent:      record.MailsSent,
		RemainingMails: record.LimitAmount - record.MailsSent,
		CreatedAt:      record.CreatedAt.String(),
		UpdatedAt:      record.UpdatedAt.String(),
	}

	return response, nil
}

func (r *MailUsageRepository) UpdateMailUsageRecord(m *model.MailUsage) error {
	var existingRecord model.MailUsage
	if err := r.DB.Where("uuid = ?", m.UUID).First(&existingRecord).Error; err != nil {
		return fmt.Errorf("error fetching mail usage record: %w", err)
	}

	existingRecord.MailsSent = m.MailsSent
	existingRecord.RemainingMails = m.RemainingMails

	if err := r.DB.Save(&existingRecord).Error; err != nil {
		return fmt.Errorf("error updating mail usage record: %w", err)
	}

	return nil
}

func (r *MailUsageRepository) GetOrCreateCurrentMailUsageRecord(subscriptionId int, limitAmount int, isPeriodDaily bool) (*model.MailUsageResponseModel, error) {
	now := time.Now()
	var periodStart, periodEnd time.Time

	if isPeriodDaily {
		periodStart = now.Truncate(24 * time.Hour)
		periodEnd = periodStart.Add(24 * time.Hour).Add(-time.Second)
	} else {
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	}

	var record model.MailUsage
	err := r.DB.Where("subscription_id = ? AND period_start = ? AND period_end = ?", subscriptionId, periodStart, periodEnd).First(&record).Error

	if err == gorm.ErrRecordNotFound {
		newRecord := &model.MailUsage{
			SubscriptionID: subscriptionId,
			PeriodStart:    periodStart,
			PeriodEnd:      periodEnd,
			LimitAmount:    limitAmount,
			MailsSent:      0,
		}
		if err := r.CreateMailUsageRecord(newRecord); err != nil {
			return nil, err
		}
		record = *newRecord
	} else if err != nil {
		return nil, err
	}

	response := &model.MailUsageResponseModel{
		ID:             int(record.ID),
		UUID:           record.UUID,
		SubscriptionID: record.SubscriptionID,
		PeriodStart:    record.PeriodStart,
		PeriodEnd:      record.PeriodEnd,
		LimitAmount:    record.LimitAmount,
		MailsSent:      record.MailsSent,
		RemainingMails: record.LimitAmount - record.MailsSent,
		CreatedAt:      record.CreatedAt.String(),
		UpdatedAt:      record.UpdatedAt.String(),
	}

	return response, nil
}
