package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type APIKeyRepository struct {
	DB *gorm.DB
}

func NewAPIkeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{
		DB: db,
	}
}

func (r *APIKeyRepository) createAPIKeyResponse(apiKey model.APIKey) model.APIKeyResponseModel {
	return model.APIKeyResponseModel{
		UUID:      apiKey.UUID,
		UserId:    apiKey.UserId,
		APIKey:    apiKey.APIKey,
		CreatedAt: apiKey.CreatedAt,
		UpdatedAt: apiKey.UpdatedAt.Format(time.RFC3339),
	}
}

func (r *APIKeyRepository) CreateAPIKey(d *model.APIKey) (*model.APIKey, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert key: %w", err)
	}
	return d, nil

}

func (r *APIKeyRepository) GetUserAPIKeyByUserId(userId string) (model.APIKeyResponseModel, error) {

	var apiKey model.APIKey

	if err := r.DB.Model(&model.APIKey{}).Where("user_id = ?", userId).First(&apiKey).Error; err != nil {
		return model.APIKeyResponseModel{}, nil
	}

	return r.createAPIKeyResponse(apiKey), nil
}

func (r *APIKeyRepository) UpdateAPIKey(d *model.APIKey) error {

	var existingAPIKey model.APIKey
	if err := r.DB.Where("user_id = ?", d.UserId).First(&existingAPIKey).Error; err != nil {
		return fmt.Errorf("failed to find API key for update: %w", err)
	}

	existingAPIKey.APIKey = d.APIKey
	existingAPIKey.UpdatedAt = time.Now()

	if err := r.DB.Save(&existingAPIKey).Error; err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	return nil
}

func (r *APIKeyRepository) CheckIfAPIKEYExists(apiKey string) (bool, error) {
	var existingAPIKey model.APIKey
	result := r.DB.Where("api_key = ?", apiKey).First(&existingAPIKey)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func (r *APIKeyRepository) DeleteAPIKey(apiKeyId string) error {
	if err := r.DB.Where("uuid = ?", apiKeyId).Delete(&model.APIKey{}).Error; err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

func (r *APIKeyRepository) FindUserWithAPIKey(apiKey string) (string, error) {
	var userID string
	if err := r.DB.Model(&model.APIKey{}).Where("api_key = ?", apiKey).Pluck("user_id", &userID).Error; err != nil {
		return "", fmt.Errorf("failed to find user with API key: %w", err)
	}
	return userID, nil
}
