package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
)

func CreateOTP(d *model.OTP) error {
	err := repository.CreateOTP(d)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveOTP(d *model.OTP) (*model.OTP, error) {
	otpData, err := repository.FindOTP(d)

	if err != nil {
		return nil, err
	}

	return otpData, err
}

func DeleteOTP(id int) error {
	err := repository.DeleteOTP(id)

	if err != nil {
		return err
	}
	return nil
}
