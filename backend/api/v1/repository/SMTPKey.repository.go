package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
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

	htime := smtpKey.UpdatedAt.String()
	response := model.SMTPDetailsResponse{
		UUID:      smtpKey.UUID,
		UserId:    smtpKey.UserId,
		KeyName:   smtpKey.KeyName,
		Password:  smtpKey.Password,
		Status:    string(smtpKey.Status),
		CreatedAt: smtpKey.CreatedAt.String(),
		UpdatedAt: &htime,
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

func (r *SMTPKeyRepository) DeleteSMTPKey(smtpKeyId string) error {
	result := r.DB.Where("uuid = ?", smtpKeyId).Delete(&model.SMTPKey{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete SMTP key: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no SMTP key found with UUID: %s", smtpKeyId)
	}
	return nil
}

func (r *SMTPKeyRepository) GetSMTPMasterKeyUserAndPass(username string, password string) (bool, error) {
	var record model.SMTPMasterKey

	if err := r.DB.Model(&model.SMTPMasterKey{}).Where("smtp_login = ? AND password = ?", username, password).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Record not found
		}
		return false, fmt.Errorf("unable to retrieve master key: %w", err) // Other errors
	}

	return true, nil // Record found
}

/*
these codes were purposely placed here for uniformity sake
they are actually meant to be in the campaing
*/

func (r *SMTPKeyRepository) MarkEmailAsDelivered(from string, to []string) error {
	// Logic to update the delivery status in the database
	for _, recipient := range to {
		if err := r.DB.Model(&model.EmailCampaignResult{}).
			Where("recipient_email = ?", recipient).
			Update("sent_at", time.Now()).
			Update("bounce_status", "").Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *SMTPKeyRepository) UpdateBounceStatus(recipientEmail, status string) error {
	return r.DB.Model(&model.EmailCampaignResult{}).
		Where("recipient_email = ?", recipientEmail).
		Update("bounce_status", status).Error
}
