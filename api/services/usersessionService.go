package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type UserSessionService struct {
	userSessionRepo *repository.UserSessionRepository
}

func NewUserSessionService(usersessionRepo *repository.UserSessionRepository) *UserSessionService {
	return &UserSessionService{
		userSessionRepo: usersessionRepo,
	}
}

func (s *UserSessionService) CreateSession(d *model.UserSessionModelStruct) (map[string]interface{}, error) {

	d.UUID = uuid.New().String()
	d.CreatedAt = time.Now()

	var response map[string]interface{}

	sessionRepo, err := s.userSessionRepo.GetSessionsByUserID(d.UserId)

	if err != nil {
		return nil, err
	}

	// If there are existing sessions, check if the new session matches any existing session
	if sessionExists := isSessionExists(sessionRepo, *d); sessionExists {
		// Return a success response
		response = map[string]interface{}{
			"message": "Session already exists for the user",
			"email":   nil,
		}
		return response, nil
	} else {
		if err := s.userSessionRepo.CreateSession(d); err != nil {
			return nil, err
		}

		response = map[string]interface{}{
			"message": "Session created successfully",
			"email":   "email sent successfully",
		}

		return response, nil
	}

}

// isSessionExists checks if the new session matches any existing session
func isSessionExists(sessionRepo []model.UserSessionResponseModel, newSession model.UserSessionModelStruct) bool {
	for _, session := range sessionRepo {
		if sessionsMatch(session, newSession) {
			return true
		}
	}
	return false
}

func sessionsMatch(sessionA model.UserSessionResponseModel, sessionB model.UserSessionModelStruct) bool {
	return ((sessionA.Device == nil && sessionB.Device == nil) || (sessionA.Device != nil && sessionB.Device != nil && *sessionA.Device == *sessionB.Device)) &&
		((sessionA.IPAddress == nil && sessionB.IPAddress == nil) || (sessionA.IPAddress != nil && sessionB.IPAddress != nil && *sessionA.IPAddress == *sessionB.IPAddress)) &&
		((sessionA.Browser == nil && sessionB.Browser == nil) || (sessionA.Browser != nil && sessionB.Browser != nil && *sessionA.Browser == *sessionB.Browser))
}

func (s *UserSessionService) GetAllSessions(userId int) ([]model.UserSessionResponseModel, error) {
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
