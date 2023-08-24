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
	err := s.otpRepository.CreateOTP(d)
	if err != nil {
		return err
	}

	return nil
}

func (s *OTPService) RetrieveOTP(d *model.OTP) (*model.OTP, error) {
	otpData, err := s.otpRepository.FindOTP(d)

	if err != nil {
		return nil, err
	}

	return otpData, err
}

func (s *OTPService) DeleteOTP(id int) error {
	err := s.otpRepository.DeleteOTP(id)

	if err != nil {
		return err
	}
	return nil
}
