package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	DB *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{
		DB: db,
	}
}

func (r *UserSessionRepository) createUserSessionResponseModel(session *model.UserSession) *model.UserSessionResponseModel {
	response := &model.UserSessionResponseModel{
		Id:        int(session.ID),
		UUID:      session.UUID,
		UserId:    session.UserId,
		Device:    session.Device,
		IPAddress: session.IPAddress,
		Browser:   session.Browser,
		CreatedAt: session.CreatedAt,
		UpdatedAt: FormatTime(session.UpdatedAt).(string),
	}

	return response
}

// Add a utility function to format time to string

func (r *UserSessionRepository) CreateSession(session *model.UserSession) error {
	if err := r.DB.Create(&session).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *UserSessionRepository) GetSessionsByUserID(userID string) ([]model.UserSessionResponseModel, error) {
	var sessions []model.UserSession

	if err := r.DB.Where("uuid = ?", userID).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	var response []model.UserSessionResponseModel
	for _, session := range sessions {
		response = append(response, *r.createUserSessionResponseModel(&session))
	}

	return response, nil
}

func (r *UserSessionRepository) DeleteSession(sessionId string) error {
	if err := r.DB.Where("uuid = ?", sessionId).Delete(&model.UserSession{}).Error; err != nil {
		return fmt.Errorf("failed to delete user session: %w", err)
	}
	return nil

}
