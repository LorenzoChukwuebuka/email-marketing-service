package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
	"strings"
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

func (s *APIKeyService) GenerateAPIKey(d *dto.APIkeyDTO) (map[string]interface{}, error) {
	uuidObj := uuid.New().String()
	apiKey := "skey-" + strings.ReplaceAll(uuidObj, "-", "")
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}
	apiKeyModel := &model.APIKey{
		UUID:   uuid.New().String(),
		UserId: d.UserId,
		APIKey: apiKey,
		Name:   d.Name,
	}

	existingKeys, err := s.APIKeyRepo.GetUserAPIKeyByUserId(d.UserId)
	if err != nil {
		return nil, err
	}

	// Check if the name already exists
	for _, existingKey := range existingKeys {
		if existingKey.Name == d.Name {
			return nil, fmt.Errorf("API key with name '%s' already exists", d.Name)
		}
	}

	// If the name doesn't exist, create the new API key
	createAPIKey, err := s.APIKeyRepo.CreateAPIKey(apiKeyModel)
	if err != nil {
		return nil, err
	}

	successMap := map[string]interface{}{
		"apiKey": createAPIKey.APIKey,
	}
	return successMap, nil
}

func (s *APIKeyService) GetAPIKey(userId string) ([]model.APIKeyResponseModel, error) {
	userApiKeys, err := s.APIKeyRepo.GetUserAPIKeyByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(userApiKeys) == 0 {
		return nil, nil
	}

	return userApiKeys, nil
}

func (s *APIKeyService) DeleteAPIKey(apiKeyId string) error {
	if err := s.APIKeyRepo.DeleteAPIKey(apiKeyId); err != nil {
		return err
	}
	return nil
}

func (s *APIKeyService) FindUserWithAPIKey(apiKey string) (string, error) {
	userId, err := s.APIKeyRepo.FindUserWithAPIKey(apiKey)

	if err != nil {
		return "", nil
	}

	return userId, nil
}
