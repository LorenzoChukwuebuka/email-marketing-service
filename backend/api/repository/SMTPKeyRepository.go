package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SMTPKeyRepository struct {
	DB *gorm.DB
}

func NewSMTPkeyRepository(db *gorm.DB) *SMTPKeyRepository {
	return &SMTPKeyRepository{
		DB: db,
	}
}

func (r *SMTPKeyRepository) createSMTPKeyResponse(smtpKey model.SMTPDetails) (*model.SMTPDetailsResponse, error) {
	response := &model.SMTPDetailsResponse{
		Id:        int(smtpKey.Id),
		UUID:      smtpKey.UUID,
		UserId:    smtpKey.UserId,
		KeyName:   smtpKey.KeyName,
		SMTPLogin: smtpKey.SMTPLogin,
		Password:  smtpKey.Password,
		CreatedAt: smtpKey.CreatedAt.String(),
	}

	if smtpKey.UpdatedAt != nil {
		response.UpdatedAt = smtpKey.UpdatedAt.Format(time.RFC3339)
	} else {
		response.UpdatedAt = ""
	}

	return response, nil
}

func (r *SMTPKeyRepository) CreateSMTPKey(d *model.SMTPDetails) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *SMTPKeyRepository) GetUserSMTPKey(userId string) (*model.SMTPDetailsResponse, error) {
	var record model.SMTPDetails

	if err := r.DB.Model(&model.SMTPDetails{}).Where("uuid = ? AND status = ?", userId, "active").First(&record).Error; err != nil {
		return nil, fmt.Errorf("unable to retrieve ")
	}

	response, err := r.createSMTPKeyResponse(record)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP key response: %w", err)
	}

	return response, nil
}

func (r *SMTPKeyRepository) UpdateSMTPKey(d *model.SMTPDetails) error {
	return nil
}
