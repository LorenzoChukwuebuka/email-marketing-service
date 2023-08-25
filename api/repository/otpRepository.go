package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"fmt"
)

type OTPRepository struct {
	DB *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{DB: db}
}

func (r *OTPRepository) CreateOTP(d *model.OTP) error {

	query := "Insert into otp (user_id,token,uuid)Values($1,$2,$3)"

	_, err := r.DB.Exec(query, d.UserId, d.Token, d.UUID)

	if err != nil {
		return err
	}

	return nil

}

func (r *OTPRepository) FindOTP(d *model.OTP) (*model.OTP, error) {

	query := "SELECT id,user_id, token, created_at,uuid FROM otp WHERE token = $1"
	row := r.DB.QueryRow(query, d.Token)

	var otp model.OTP
	err := row.Scan(&otp.Id, &otp.UserId, &otp.Token, &otp.CreatedAt, &otp.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("otp does not exist: %w", err) // OTP not found, return nil without an error
		}
		return nil, err
	}

	return &otp, nil
}

func (r *OTPRepository) DeleteOTP(id int) error {

	query := "DELETE FROM otp WHERE id = $1"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
