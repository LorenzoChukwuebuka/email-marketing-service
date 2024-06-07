package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
)

type OTPService struct {
	otpRepository *repository.OTPRepository
}

func NewOTPService(otpRepo *repository.OTPRepository) *OTPService {
	return &OTPService{
		otpRepository: otpRepo,
	}
}

func (s *OTPService) CreateOTP(d *model.OTP) error {
	if err := s.otpRepository.CreateOTP(d); err != nil {
		return err
	}

	return nil
}

func (s *OTPService) RetrieveOTP(d *model.OTP) (*model.OTP, error) {
	otpData, err := s.otpRepository.FindOTP(d)

	if err != nil {
		return nil, err
	}

	return otpData, nil
}

func (s *OTPService) DeleteOTP(id int) error {
	if err := s.otpRepository.DeleteOTP(id); err != nil {
		return err
	}

	return nil
}
