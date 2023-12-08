package services

import (
	"fmt"

	"github.com/google/uuid"
)

type ApiKeyService struct {
}

func NewAPIKeyService() *ApiKeyService {
	return &ApiKeyService{}
}

func (s *ApiKeyService) GenerateAPIKey(userId int) (string, error) {

	uuidObj := uuid.New()

	fmt.Println(uuidObj)

	return "", nil

}
