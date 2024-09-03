package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type DomainRepository struct {
	DB *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{
		DB: db,
	}
}

func (r *DomainRepository) CreateDomain(d *model.Domains) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert domain: %w", err)
	}
	return nil
}

func (r *DomainRepository) CheckIfDomainExists(d *model.Domains) (bool, error) {
	result := r.DB.Where("domain = ? AND user_id =?", d.Domain, d.UserID).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
