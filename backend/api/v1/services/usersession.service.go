package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
	"log"
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

	session := &model.UserSession{UUID: uuid.New().String(), UserId: d.UserId, Device: d.Device, IPAddress: d.IPAddress, Browser: d.Browser}

	var response map[string]interface{}

	sessionRepo, err := s.userSessionRepo.GetSessionsByUserID(d.UserId)

	if err != nil {
		return nil, err
	}

	// If there are existing sessions, check if the new session matches any existing session
	if sessionExists := isSessionExists(sessionRepo, *session); sessionExists {
		// Return a success response
		response = map[string]interface{}{
			"message": "Session already exists for the user",
			"email":   nil,
		}
		return response, nil
	} else {
		if err := s.userSessionRepo.CreateSession(session); err != nil {
			return nil, err
		}

		//this will be in a separate go routine abeg I no get strength

		resultChan := make(chan Result, 1)

		defer close(resultChan)

		go s.sendDeviceVerificationMail(session, resultChan)

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
		UUID: d.UserId,
	}
	userRepo, err := s.userRepo.FindUserById(userStruct)

	if err != nil {
		resultChan <- Result{Error: fmt.Errorf("failed to find user by ID: %w", err)}
		return
	}
	userEmail := userRepo.Email
	userName := userRepo.FullName
	code := utils.GenerateOTP(8)
	err = mailer.DeviceVerificationMail(userName, userEmail, d, code)

	if err != nil {
		resultChan <- Result{Error: fmt.Errorf("email sending failed: %w", err)}
		return
	}

	resultChan <- Result{Success: true}
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
