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

var EmptyAPIKeyResponse = model.APIKeyResponseModel{}

func NewAPIKeyService(apiRepo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		APIKeyRepo: apiRepo,
	}
}

func (s *APIKeyService) GenerateAPIKey(userId int) (map[string]interface{}, error) {
	uuidObj := uuid.New().String()
	apiKey := "skey-" + uuidObj

	apiKeyModel := &model.APIKey{
		UUID:      uuid.New().String(),
		UserId:    userId,
		APIKey:    apiKey,
		CreatedAt: time.Now(),
	}

	existingAPIKey, err := s.APIKeyRepo.GetUserAPIKeyByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("error checking existing API key: %v", err)
	}

	if existingAPIKey != EmptyAPIKeyResponse {
		// If an existing API key is found, update it

		err := s.APIKeyRepo.UpdateAPIKey(apiKeyModel)
		if err != nil {
			return nil, fmt.Errorf("error updating API key: %v", err)
		}
	} else {
		// If no existing API key, create a new one
		_, err := s.APIKeyRepo.CreateAPIKey(apiKeyModel)
		if err != nil {
			return nil, fmt.Errorf("error generating API key: %v", err)
		}
	}

	successMap := map[string]interface{}{
		"apiKey": apiKeyModel.APIKey,
	}

	return successMap, nil
}

func (s *APIKeyService) GetAPIKey(userId int) (model.APIKeyResponseModel, error) {
	userApiKey, err := s.APIKeyRepo.GetUserAPIKeyByUserId(userId)

	if err != nil {
		return EmptyAPIKeyResponse, err
	}

	if userApiKey == EmptyAPIKeyResponse {
		return EmptyAPIKeyResponse, fmt.Errorf("no api key generated yet")
	}
	return userApiKey, nil
}

func (s *APIKeyService) DeleteAPIKey(apiKeyId string) error {
	if err := s.APIKeyRepo.DeleteAPIKey(apiKeyId); err != nil {
		return err
	}
	return nil
}

func (s *APIKeyService) FindUserWithAPIKey(apiKey string) (int, error) {
	userId, err := s.APIKeyRepo.FindUserWithAPIKey(apiKey)

	if err != nil {
		return 0, nil
	}

	return userId, nil
}
