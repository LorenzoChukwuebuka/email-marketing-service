package repository

import (
	"email-marketing-service/api/v1/model"
	"gorm.io/gorm"
)

type SupportRepository struct {
	DB *gorm.DB
}

func NewSupportRepository(db *gorm.DB) *SupportRepository {
	return &SupportRepository{DB: db}
}

func (r *SupportRepository) CreateSupportTicket(d *model.SupportTicket) (uint, error) {
	if err := r.DB.Create(d).Error; err != nil {
		return 0, err
	}
	return d.ID, nil
}

func (r *SupportRepository) CreateSupportTicketFile(d *model.TicketFiles) error {
	if err := r.DB.Create(d).Error; err != nil {
		return err
	}

	return nil
}
