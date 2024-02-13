package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
)

type OTPRepository struct {
	DB *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{DB: db}
}

func (r *OTPRepository) CreateOTP(d *model.OTP) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert otp: %w", err)
	}
	return nil
}

func (r *OTPRepository) FindOTP(d *model.OTP) (*model.OTP, error) {

	result := r.DB.Where("token = ?", d.Token).First(&d)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("token not found:%w", result.Error)
		}
		return nil, result.Error
	}

	return d, nil
}

func (r *OTPRepository) DeleteOTP(id int) error {
	var otp *model.OTP

	// Fetch the OTP record from the database
	if err := r.DB.First(&otp, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("OTP not found")
		} else {
			fmt.Printf("Error querying database: %v\n", err)
		}
		return nil
	}

	// Delete the OTP record
	if err := r.DB.Delete(&otp).Error; err != nil {
		fmt.Printf("Error deleting OTP: %v\n", err)
		return nil
	}
	return nil
}
