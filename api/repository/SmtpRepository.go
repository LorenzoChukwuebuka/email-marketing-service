package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
)

type SMTPWebHookRepository struct {
	DB gorm.DB
}

func  NewSMTPWebHookRepository(db *gorm.DB) *SMTPWebHookRepository {
	return &SMTPWebHookRepository{
		DB: *db,
	}
}

func (r *SMTPWebHookRepository) CreateReport(d *model.SentEmails) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert sent mail: %w", err)
	}
	return nil
}


func (r *SMTPWebHookRepository) UpdateReport(d *model.SentEmails) error {
	return nil 
}
