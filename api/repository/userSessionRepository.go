package repository

import (
	"email-marketing-service/api/model"
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

func (r *UserSessionRepository) CreateSession(session *model.UserSession) error {
	if err := r.DB.Create(&session).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *UserSessionRepository) GetSessionsByUserID(userID int) ([]model.UserSessionResponseModel, error) {
	return nil, nil
}

func (r *UserSessionRepository) DeleteSession(sessionId string) error {

	return nil

}
