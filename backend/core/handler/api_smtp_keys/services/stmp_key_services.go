package services

import (
	"context"

	"email-marketing-service/core/handler/api_smtp_keys/dto"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/google/uuid"
)

var (
	cfg = config.LoadEnv()
)

type SMTPKeyService struct {
	store db.Store
}

func NewSMTPKeyService(store db.Store) *SMTPKeyService {
	return &SMTPKeyService{store: store}
}

func (s *SMTPKeyService) GenerateNewSMTPMasterPassword(ctx context.Context, userid string) (any, error) {
	// Generate new password and SMTP login
	password := helper.GenerateOTP(15)
	smtpLogin := helper.GenerateOTP(8) + "@" + cfg.SMTP_SERVER

	userID, err := uuid.Parse(userid)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	// Update the master key with the new password and login
	err = s.store.UpdateSMTPKeyMasterPasswordAndLogin(ctx, db.UpdateSMTPKeyMasterPasswordAndLoginParams{
		Password:  password,
		SmtpLogin: smtpLogin,
		UserID:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update SMTP master key: %w", err)
	}

	// Fetch all SMTP keys for the user
	smtpKeys, err := s.store.GetUserSMTPKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user SMTP keys: %w", err)
	}

	//Update the SMTP login for each SMTP key associated with the user

	for _, smtpKey := range smtpKeys {
		err := s.store.UpdateSMTPKeyLogin(ctx, db.UpdateSMTPKeyLoginParams{
			UserID:    userID,
			SmtpLogin: smtpLogin,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update SMTP key login for key '%s': %w", smtpKey.KeyName, err)
		}
	}

	// Return the new password and login
	result := map[string]interface{}{
		"password": password,
		"login":    smtpLogin,
	}

	return result, nil
}

func (s *SMTPKeyService) GetSMTPKeys(ctx context.Context, userId string) (any, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	masterSMTP, err := s.store.GetMasterSMTPKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	keys, err := s.store.GetUserSmtpKeys(ctx, userID)
	if err != nil {
		return nil, err
	}

	var formattedKeys []map[string]interface{}
	for _, key := range keys {
		formattedKey := map[string]interface{}{
			"id":         key.ID,
			"company_id": key.CompanyID,
			"user_id":    key.UserID,
			"key_name":   key.KeyName,
			"password":   key.Password,
			"status":     key.Status,
			"smtp_login": key.SmtpLogin,
			"created_at": key.CreatedAt,
			"updated_at": key.UpdatedAt,
		}

		if key.DeletedAt.Valid {
			formattedKey["deleted_at"] = key.DeletedAt.Time.Format("2006-01-02")
		} else {
			formattedKey["deleted_at"] = nil
		}

		formattedKeys = append(formattedKeys, formattedKey)
	}

	result := map[string]interface{}{
		"smtp_master_password": masterSMTP.Password,
		"smtp_master":          masterSMTP.KeyName,
		"smtp_login":           masterSMTP.SmtpLogin,
		"smtp_port":            cfg.SMTP_PORT,
		"smtp_server":          cfg.SMTP_SERVER,
		"smtp_master_status":   masterSMTP.Status,
		"smtp_created_at":      masterSMTP.CreatedAt,
		"keys":                 formattedKeys,
	}

	return result, nil
}

func (s *SMTPKeyService) CreateSMTPKey(ctx context.Context, d *dto.SMTPKeyRequestDTO) (map[string]interface{}, error) {
	if err := helper.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	userID, err := uuid.Parse(d.UserId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	companyId, err := uuid.Parse(d.CompanyID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	// Fetch the SMTP login from the master key
	masterSMTP, err := s.store.GetMasterSMTPKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve SMTP master key: %w", err)
	}

	// Generate the password for the SMTP key
	password := helper.GenerateOTP(15)

	// Check if the key name already exists
	keyNameExist, err := s.store.GetUserSmtpKeys(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, key := range keyNameExist {
		if key.KeyName == d.KeyName {
			return nil, fmt.Errorf("SMTP key with name '%s' already exists for this user", d.KeyName)
		}
	}

	// Save the new SMTP key
	_, err = s.store.CreateSMTPKey(ctx, db.CreateSMTPKeyParams{
		CompanyID: companyId,
		UserID:    userID,
		KeyName:   d.KeyName,
		Password:  password,
		Status:    string(enums.KeyActive),
		SmtpLogin: masterSMTP.SmtpLogin,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"password": password,
	}, nil
}

func (s *SMTPKeyService) ToggleSMTPKeyStatus(ctx context.Context, userId string, smtpKeyId string) error {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	smtpkeyID, err := uuid.Parse(smtpKeyId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	// Retrieve the SMTP key record
	smtpKey, err := s.store.GetSMTPKeyByID(ctx, db.GetSMTPKeyByIDParams{
		UserID: userID,
		ID:     smtpkeyID,
	})
	if err != nil {
		return common.ErrRecordNotFound
	}

	// Toggle the status
	var newStatus string
	if smtpKey.Status == string(enums.KeyActive) {
		newStatus = string(enums.KeyInactive)
	} else {
		newStatus = string(enums.KeyActive)
	}

	// Update the SMTP key status
	err = s.store.UpdateSMTPKeyStatus(ctx, db.UpdateSMTPKeyStatusParams{
		Status: newStatus,
		ID:     smtpkeyID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("unable to update SMTP key status: %w", err)
	}

	return nil
}

func (s *SMTPKeyService) DeleteSMTPKey(ctx context.Context, smtpkeyId string) error {
	smtpkeyID, err := uuid.Parse(smtpkeyId)
	if err != nil {
		return common.ErrInvalidUUID
	}
	if err := s.store.SoftDeleteSMTPKey(ctx, smtpkeyID); err != nil {
		return err
	}
	return nil
}
