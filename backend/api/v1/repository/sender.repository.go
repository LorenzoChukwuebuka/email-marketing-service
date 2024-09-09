package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type SenderRepository struct {
	DB *gorm.DB
}

func NewSenderRepository(db *gorm.DB) *SenderRepository {
	return &SenderRepository{
		DB: db,
	}
}

func (r *SenderRepository) CreateSender(d *model.Sender) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert domain: %w", err)
	}
	return nil
}

func (r *SenderRepository) CheckIfSenderExists(d *model.Sender) (bool, error) {
	result := r.DB.Where("email = ? AND name = ? AND user_id =?", d.Email, d.Name, d.UserID).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
