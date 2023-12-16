package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
)

type UserSessionService struct {
	userSessionRepo *repository.UserSessionRepository
}

func NewUserSessionService(usersessionRepo *repository.UserSessionRepository) *UserSessionService {
	return &UserSessionService{
		userSessionRepo: usersessionRepo,
	}
}

func (s *UserSessionService) CreateSession(d *model.UserSessionModelStruct) (*model.UserSessionModelStruct, error) {

	return nil, nil
}
