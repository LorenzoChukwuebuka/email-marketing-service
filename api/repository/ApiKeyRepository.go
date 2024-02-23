package repository

import (
	"email-marketing-service/api/model"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	DB *gorm.DB
}

func NewAPIkeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{
		DB: db,
	}
}

func (r *APIKeyRepository) CreateAPIKey(d *model.APIKey) (*model.APIKey, error) {
 return nil,nil

}

func (r *APIKeyRepository) GetUserAPIKeyByUserId(userId int) (*model.APIKeyResponseModel, error) {
	 return nil,nil
}

func (r *APIKeyRepository) UpdateAPIKey(d *model.APIKey) error {
	 
	return nil
}

func (r *APIKeyRepository) CheckIfAPIKEYExists(apiKey string) (bool, error) {

 return false,nil
}

func (r *APIKeyRepository) DeleteAPIKey(apiKeyId string) error {
	 

	return nil
}


func (r *APIKeyRepository) FindUserWithAPIKey(apiKey string) (int, error) {
    return 0, nil
}
