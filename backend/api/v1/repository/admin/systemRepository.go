package adminrepository

import (
	adminmodel "email-marketing-service/api/v1/model/admin"

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
