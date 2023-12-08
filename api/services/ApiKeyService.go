package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type APIKeyService struct {
	APIKeyRepo *repository.APIKeyRepository
}

func NewAPIKeyService(apiRepo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		APIKeyRepo: apiRepo,
	}
}

func (s *APIKeyService) GenerateAPIKey(userId int) (map[string]interface{}, error) {
	uuidObj := uuid.New().String()
	// Concatenate strings
	apiKey := "skey-" + uuidObj

	apiKeyModel := &model.APIKeyModel{
		UUID:      uuid.New().String(),
		UserId:    userId,
		APIKey:    apiKey,
		CreatedAt: time.Now(),
	}

	apiRepo, err := s.APIKeyRepo.CreateAPIKey(apiKeyModel)

	if err != nil {
		return nil, fmt.Errorf("error generating API key: %v", err)
	}

	successMap := map[string]interface{}{
		"apiKey": apiRepo.APIKey,
	}

	return successMap, nil
}

func (s *APIKeyService) UpdateAPIKey(userId int) (map[string]interface{}, error) {
	return nil, nil
}

func (s *APIKeyService) GetAPIKey(userId int) (*model.APIKeyResponseModel, error) {

	return nil, nil
}
