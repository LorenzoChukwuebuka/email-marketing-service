package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SMTPKeyRepository struct {
	DB *gorm.DB
}

func NewSMTPkeyRepository(db *gorm.DB) *SMTPKeyRepository {
	return &SMTPKeyRepository{
		DB: db,
	}
}

func (r *SMTPKeyRepository) createSMTPKeyResponse(smtpKey model.SMTPKey) (model.SMTPDetailsResponse, error) {
	response := model.SMTPDetailsResponse{
		UUID:      smtpKey.UUID,
		UserId:    smtpKey.UserId,
		KeyName:   smtpKey.KeyName,
		Password:  smtpKey.Password,
		Status:    string(smtpKey.Status),
		CreatedAt: FormatTime(smtpKey.CreatedAt).(string),
		UpdatedAt: FormatTime(smtpKey.UpdatedAt).(*string),
	}

	if smtpKey.DeletedAt.Valid {
		formatted := smtpKey.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &formatted
	}

	return response, nil
}

func (r *SMTPKeyRepository) CreateSMTPKey(d *model.SMTPKey) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to create key : %w", err)
	}
	return nil
}

func (r *SMTPKeyRepository) GetUserSMTPKey(userId string) ([]model.SMTPDetailsResponse, error) {
	var records []model.SMTPKey

	if err := r.DB.Where("user_id = ?", userId).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("unable to retrieve SMTP keys: %w", err)
	}

	var response []model.SMTPDetailsResponse

	for _, smtpKey := range records {
		detailsResponse, err := r.createSMTPKeyResponse(smtpKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create SMTP key response: %w", err)
		}
		response = append(response, detailsResponse)
	}

	return response, nil
}

func (r *SMTPKeyRepository) CreateSMTPMasterKey(d *model.SMTPMasterKey) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to create master key : %w", err)
	}
	return nil
}

func (r *SMTPKeyRepository) UpdateSMTPKeyMasterPasswordAndLogin(userId string, password string, smtpLogin string) error {
	err := r.DB.Model(&model.SMTPMasterKey{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"password":   password,
		"smtp_login": smtpLogin,
		"updated_at": time.Now().UTC(),
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update smtp master key: %w", err)
	}

	return nil
}

func (r *SMTPKeyRepository) GetSMTPMasterKey(userId string) (*model.SMTPMasterKey, error) {

	var record model.SMTPMasterKey

	if err := r.DB.Model(&model.SMTPMasterKey{}).Where("user_id = ? ", userId).First(&record).Error; err != nil {
		return nil, fmt.Errorf("unable to retrieve master key")
	}

	return &record, nil
}

func (r *SMTPKeyRepository) ToggleSMTPKeyStatus(userId string, smtpkeyId string) error {
	// Define a variable to hold the SMTP key record
	var smtpKey model.SMTPKey

	// Retrieve the SMTP key record by userId and smtpkeyId
	if err := r.DB.Where("user_id = ? AND uuid = ?", userId, smtpkeyId).First(&smtpKey).Error; err != nil {
		return fmt.Errorf("unable to retrieve SMTP key: %w", err)
	}

	// Toggle the status
	if smtpKey.Status == model.KeyActive {
		smtpKey.Status = model.KeyStatus(model.KeyInactive)
	} else {
		smtpKey.Status = model.KeyActive
	}

	currentTime := time.Now().UTC()
	smtpKey.UpdatedAt = currentTime

	// Save the updated record back to the database
	if err := r.DB.Save(&smtpKey).Error; err != nil {
		return fmt.Errorf("unable to update SMTP key status: %w", err)
	}

	return nil
}

func (r *SMTPKeyRepository) DeleteSMTPKey(smtkeyId string) error {
	if err := r.DB.Where("uuid = ?", smtkeyId).Delete(&model.APIKey{}).Error; err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}
