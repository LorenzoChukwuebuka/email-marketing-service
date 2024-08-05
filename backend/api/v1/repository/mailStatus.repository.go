package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type MailStatusRepository struct {
	DB gorm.DB
}

func NewMailStatusRepository(db *gorm.DB) *MailStatusRepository {
	return &MailStatusRepository{
		DB: *db,
	}
}

func (r *MailStatusRepository) CreateStatus(d *model.SentEmails) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert sent mail: %w", err)
	}
	return nil
}

func (r *MailStatusRepository) UpdateReport(d *model.SentEmails) error {
	err := r.DB.Model(&model.SentEmails{}).Where("uuid = ?", d.UUID).Updates(map[string]interface{}{
		"status":     d.Status,
		"updated_at": d.UpdatedAt,
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update sent mail: %w", err)
	}

	return nil
}

func (r *MailStatusRepository) GetSentEmailByRecipient(email string) (*model.SentEmails, error) {
	var sentEmail model.SentEmails
	err := r.DB.Where("recipient = ?", email).Order("created_at DESC").Limit(1).Last(&sentEmail).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if no record is found
		}
		return nil, fmt.Errorf("failed to get sent email by recipient: %w", err)
	}
	return &sentEmail, nil
}


