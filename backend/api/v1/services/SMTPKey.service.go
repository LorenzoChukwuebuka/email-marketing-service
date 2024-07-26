package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
)

type SMTPKeyService struct {
	SMTPKeyRepo *repository.SMTPKeyRepository
}

func NewSMTPKeyService(smtpKeyRepo *repository.SMTPKeyRepository) *SMTPKeyService {
	return &SMTPKeyService{
		SMTPKeyRepo: smtpKeyRepo,
	}
}

var (
	smtplog = config.SMTP_SERVER
)

func (s *SMTPKeyService) GenerateNewSMTPMasterPassword(userid string) (map[string]interface{}, error) {

	password := utils.GenerateOTP(15)
	smtpLogin := utils.GenerateOTP(8) + "@" + smtplog

	err := s.SMTPKeyRepo.UpdateSMTPKeyMasterPasswordAndLogin(userid, password, smtpLogin)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"password": password,
		"login":    smtpLogin,
	}

	return result, nil
}

func (s *SMTPKeyService) GetSMTPKeys(userId string) (map[string]interface{}, error) {

	masterSMTP, err := s.SMTPKeyRepo.GetSMTPMasterKey(userId)
	if err != nil {
		return nil, err
	}

	Keys, err := s.SMTPKeyRepo.GetUserSMTPKey(userId)

	if err != nil {
		return nil, err
	}

	var keyValue interface{}

	if len(Keys) == 0 {
		keyValue = []interface{}{} // Empty array
	} else {
		keyValue = Keys
	}

	result := map[string]interface{}{
		"smtp_master_password": masterSMTP.Password,
		"smtp_master":          masterSMTP.KeyName,
		"smtp_login":           masterSMTP.SMTPLogin,
		"smtp_port":            config.SMTP_PORT,
		"smtp_server":          config.SMTP_SERVER,
		"smtp_master_status":   masterSMTP.Status,
		"smtp_created_at":      masterSMTP.CreatedAt,
		"keys":                 keyValue,
	}

	return result, nil
}

func (s *SMTPKeyService) CreateSMTPKey(d *dto.SMTPKeyDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	password := utils.GenerateOTP(15)

	keyNameExist, err := s.SMTPKeyRepo.GetUserSMTPKey(d.UserId)

	if err != nil {
		return nil, err
	}

	for _, key := range keyNameExist {
		if key.KeyName == d.KeyName {
			return nil, fmt.Errorf("SMTP key with name '%s' already exists for this user", d.KeyName)
		}
	}

	smtpKey := &model.SMTPKey{
		UUID:     uuid.New().String(),
		UserId:   d.UserId,
		KeyName:  d.KeyName,
		Password: password,
		Status:   model.KeyActive,
	}

	err = s.SMTPKeyRepo.CreateSMTPKey(smtpKey)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"password": password,
	}, nil

}

func (s *SMTPKeyService) ToggleSMTPKeyStatus(userId string, smtpKeyId string) error {
	err := s.SMTPKeyRepo.ToggleSMTPKeyStatus(userId, smtpKeyId)

	if err != nil {
		return err
	}

	return nil

}

func (s *SMTPKeyService) DeleteSMTPKey(smtpkeyId string) error {

	if err := s.SMTPKeyRepo.DeleteSMTPKey(smtpkeyId); err != nil {
		return err
	}

	return nil
}
