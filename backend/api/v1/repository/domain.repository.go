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

func (r *DomainRepository) GetDomain(uuid string) (*model.DomainsResponse, error) {
	var domain model.Domains
	if err := r.DB.Where("uuid = ?", uuid).First(&domain).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("domain not found")
		}
		return nil, err
	}

	response := &model.DomainsResponse{
		UUID:          domain.UUID,
		UserID:        domain.UserID,
		Domain:        domain.Domain,
		TXTRecord:     domain.TXTRecord,
		DMARCRecord:   domain.DMARCRecord,
		DKIMSelector:  domain.DKIMSelector,
		DKIMPublicKey: domain.DKIMPublicKey,
		Verified:      domain.Verified,
		CreatedAt:     domain.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     domain.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if domain.DeletedAt.Valid {
		deletedAt := domain.DeletedAt.Time.Format("2006-01-02 15:04:05")
		response.DeletedAt = &deletedAt
	}

	return response, nil
}

func (r *DomainRepository) UpdateDomain(d *model.Domains) error {
	return r.DB.Model(&model.Domains{}).Where("uuid = ?", d.UUID).Updates(d).Error
}

func (r *DomainRepository) DeleteDomain(id string) error {
	if err := r.DB.Where("uuid = ?", id).Delete(&model.Domains{}).Error; err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}

func (r *DomainRepository) GetAllDomains(userID string) (*[]model.DomainsResponse, error) {
	var domains []model.Domains
	if err := r.DB.Where("user_id = ?", userID).Find(&domains).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch domains: %w", err)
	}

	var responses []model.DomainsResponse
	for _, domain := range domains {
		response := model.DomainsResponse{
			UUID:          domain.UUID,
			UserID:        domain.UserID,
			Domain:        domain.Domain,
			TXTRecord:     domain.TXTRecord,
			DMARCRecord:   domain.DMARCRecord,
			DKIMSelector:  domain.DKIMSelector,
			DKIMPublicKey: domain.DKIMPublicKey,
			Verified:      domain.Verified,
			CreatedAt:     domain.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     domain.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if domain.DeletedAt.Valid {
			deletedAt := domain.DeletedAt.Time.Format("2006-01-02 15:04:05")
			response.DeletedAt = &deletedAt
		}

		responses = append(responses, response)
	}

	return &responses, nil
}
