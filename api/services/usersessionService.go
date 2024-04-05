package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
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

func (s *UserSessionService) CreateSession(d *model.UserSession) (map[string]interface{}, error) {

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

		//this will be in a separate go routine abeg I no get strength

		resultChan := make(chan Result, 1)

		defer close(resultChan)

		go s.sendDeviceVerificationMail(d, resultChan)

		result := <-resultChan

		if result.Success {
			log.Println("Mail Result:", result.Message)
		} else {
			log.Println("Mail Error:", result.Error)
		}

		response = map[string]interface{}{
			"message": "Session created successfully",
			"email":   "email sent successfully",
		}

		return response, nil
	}

}

// isSessionExists checks if the new session matches any existing session
func isSessionExists(sessionRepo []model.UserSessionResponseModel, newSession model.UserSession) bool {
	for _, session := range sessionRepo {
		if sessionsMatch(session, newSession) {
			return true
		}
	}
	return false
}

func sessionsMatch(sessionA model.UserSessionResponseModel, sessionB model.UserSession) bool {
	return ((sessionA.Device == nil && sessionB.Device == nil) || (sessionA.Device != nil && sessionB.Device != nil && *sessionA.Device == *sessionB.Device)) &&
		((sessionA.IPAddress == nil && sessionB.IPAddress == nil) || (sessionA.IPAddress != nil && sessionB.IPAddress != nil && *sessionA.IPAddress == *sessionB.IPAddress)) &&
		((sessionA.Browser == nil && sessionB.Browser == nil) || (sessionA.Browser != nil && sessionB.Browser != nil && *sessionA.Browser == *sessionB.Browser))
}

func (s *UserSessionService) sendDeviceVerificationMail(d *model.UserSession, resultChan chan Result) {
	userStruct := &model.User{
		ID: d.UserId,
	}
	userRepo, err := s.userRepo.FindUserById(userStruct)

	if err != nil {
		resultChan <- Result{Error: fmt.Errorf("failed to find user by ID: %w", err)}
		return
	}
	userEmail := userRepo.Email
	userName := userRepo.FullName
	code := utils.GenerateOTP(8)
	err = mail.DeviceVerificationMail(userName, userEmail, d, code)

	if err != nil {
		resultChan <- Result{Error: fmt.Errorf("email sending failed: %w", err)}
		return
	}

	resultChan <- Result{Success: true}
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
