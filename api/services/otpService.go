package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
)

type OTPService struct{}

func NewRepository() *repository.OTPRepository {
	return &repository.OTPRepository{}
}

func (s *OTPService) CreateOTP(d *model.OTP) error {
	err := NewRepository().CreateOTP(d)
	if err != nil {
		return err
	}

	return nil
}

func (s *OTPService) RetrieveOTP(d *model.OTP) (*model.OTP, error) {
	otpData, err := NewRepository().FindOTP(d)

	if err != nil {
		return nil, err
	}

	return otpData, err
}

func (s *OTPService) DeleteOTP(id int) error {
	err := NewRepository().DeleteOTP(id)

	if err != nil {
		return err
	}
	return nil
}
