package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type DomainRepository struct {
	DB *gorm.DB
}

var ErrDomainNotFound = errors.New("domain not found")

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
		UUID:           domain.UUID,
		UserID:         domain.UserID,
		Domain:         domain.Domain,
		TXTRecord:      domain.TXTRecord,
		DMARCRecord:    domain.DMARCRecord,
		DKIMSelector:   domain.DKIMSelector,
		DKIMPublicKey:  domain.DKIMPublicKey,
		DKIMPrivateKey: domain.DKIMPrivateKey,
		MXRecord:       domain.MXRecord,
		Verified:       domain.Verified,
		CreatedAt:      domain.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      domain.UpdatedAt.Format("2006-01-02 15:04:05"),
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

func (r *DomainRepository) GetAllDomains(userID string, searchQuery string, params PaginationParams) (PaginatedResult, error) {
	var domains []model.Domains
	query := r.DB.Model(&model.Domains{}).Where("user_id = ?", userID)

	if searchQuery != "" {
		query = query.Where("domain ILIKE ?", "%"+searchQuery+"%")
	}

	query.Order("created_at DESC")

	paginatedResult, err := Paginate(query, params, &domains)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
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
			SPFRecord:     domain.SPFRecord,
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

	paginatedResult.Data = responses

	return paginatedResult, nil
}

func (r *DomainRepository) FindDomain(userId string, domainName string) (*model.DomainsResponse, error) {
	var domain model.Domains
	if err := r.DB.Where("user_id = ? AND domain = ?", userId, domainName).First(&domain).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDomainNotFound
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
		SPFRecord:     domain.SPFRecord,
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
