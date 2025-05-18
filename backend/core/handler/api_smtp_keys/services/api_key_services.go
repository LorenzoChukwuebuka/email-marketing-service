package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/api_smtp_keys/dto"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type APIKeyService struct {
	store db.Store
}

func NewAPIKeyService(store db.Store) *APIKeyService {
	return &APIKeyService{store: store}
}

func (s *APIKeyService) GenerateAPIKey(ctx context.Context, d *dto.APIkeyRequestDTO) (any, error) {
	uuidObj := uuid.New().String()
	apiKey := "skey-" + strings.ReplaceAll(uuidObj, "-", "")

	if err := helper.ValidateData(d); err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(d.UserId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	companyId, err := uuid.Parse(d.CompanyID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	existingKeys, err := s.store.GetAPIKeysByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if the name already exists
	for _, existingKey := range existingKeys {
		if existingKey.Name == d.Name {
			return nil, fmt.Errorf("api key name already exists")
		}
	}

	// If the name doesn't exist, create the new API key
	createAPIKey, err := s.store.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		CompanyID: companyId,
		Name:      d.Name,
		ApiKey:    apiKey,
	})

	if err != nil {
		return nil, err
	}

	successMap := map[string]interface{}{
		"apiKey": createAPIKey.ApiKey,
	}
	return successMap, nil
}

func (s *APIKeyService) GetAPIKey(ctx context.Context, userId string) (any, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	userApiKeys, err := s.store.GetAPIKeysByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(userApiKeys) == 0 {
		return nil, nil
	}

	return userApiKeys, nil
}

func (s *APIKeyService) DeleteAPIKey(ctx context.Context, apiKeyId uuid.UUID) error {
	if err := s.store.DeleteAPIKey(ctx, apiKeyId); err != nil {
		return err
	}
	return nil
}

func (s *APIKeyService) FindUserWithAPIKey(ctx context.Context, apiKey string) (string, error) {
	userId, err := s.store.FindUserWithAPIKey(ctx, apiKey)
	if err != nil {
		return "", nil
	}
	return userId.UserID.String(), nil
}
