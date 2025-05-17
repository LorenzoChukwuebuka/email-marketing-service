package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"
)

type UserSessionService struct {
	userSessionRepo *repository.UserSessionRepository
	userRepo        *repository.UserRepository
}

type Result struct {
	Success bool
	Message string
	Error   error
}

func NewUserSessionService(usersessionRepo *repository.UserSessionRepository, userRepository *repository.UserRepository) *UserSessionService {
	return &UserSessionService{
		userSessionRepo: usersessionRepo,
		userRepo:        userRepository,
	}
}

func (s *UserSessionService) CreateSession(d *dto.UserSession) (map[string]interface{}, error) {
	return nil, nil
}

func (s *UserSessionService) GetAllSessions(userId string) ([]model.UserSessionResponseModel, error) {
	sessionRepo, err := s.userSessionRepo.GetSessionsByUserID(userId)

	if err != nil {
		return nil, err
	}

	if len(sessionRepo) == 0 {
		return nil, fmt.Errorf("no records found")
	}

	return sessionRepo, nil
}

func (s *UserSessionService) DeleteSession(sessionId string) error {

	err := s.userSessionRepo.DeleteSession(sessionId)

	if err != nil {
		return err
	}

	return nil
}
