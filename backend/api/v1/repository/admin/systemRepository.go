package adminrepository

import (
	adminmodel "email-marketing-service/api/v1/model/admin"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SystemRepository struct {
	DB *gorm.DB
}

func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{
		DB: db,
	}
}

func (r *SystemRepository) CreateSMTPSettings(settings *adminmodel.SystemsSMTPSetting) error {
	return r.DB.Create(settings).Error
}

func (r *SystemRepository) DomainExists(domain string) (bool, error) {
	var exists bool
	err := r.DB.Model(&adminmodel.SystemsSMTPSetting{}).
		Select("count(*) > 0").
		Where("domain = ?", domain).
		Find(&exists).Error

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SystemRepository) GetSMTPSettings(domain string) (*adminmodel.SystemsSMTPSetting, error) {
	var settings adminmodel.SystemsSMTPSetting
	err := r.DB.Where("domain = ?", domain).First(&settings).Error
	if err != nil {
		return nil, err

	}
	return &settings, nil
}

func (r *SystemRepository) UpdateSMTPSettings(settings *adminmodel.SystemsSMTPSetting) error {
	return r.DB.Save(settings).Error
}

func (r *SystemRepository) DeleteSettings(domain string) error {


	fmt.Printf("%s",domain)

	
	// Use Unscoped() to find the record even if it's soft deleted
	var setting adminmodel.SystemsSMTPSetting
	result := r.DB.Where("domain = ?", domain).First(&setting)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("domain not found: %s", domain)
		}
		return fmt.Errorf("failed to find domain: %w", result.Error)
	}

	// Use Delete() without Unscoped() to perform soft delete
	result = r.DB.Delete(&setting)
	if result.Error != nil {
		return fmt.Errorf("failed to delete domain: %w", result.Error)
	}

	return nil
}
