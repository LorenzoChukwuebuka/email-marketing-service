package repository

import (
	"database/sql"
	"email-marketing-service/api/database"
	"email-marketing-service/api/model"
)

type OTPRepository struct{}

func (r *OTPRepository) CreateOTP(d *model.OTP) error {

	db, err := database.InitDB()
	if err != nil {
		return nil
	}
	defer db.Close()

	query := "Insert into otp (user_id,token,uuid)Values($1,$2,$3)"

	_, err = db.Exec(query, d.UserId, d.Token, d.UUID)

	if err != nil {
		return err
	}

	return nil

}

func (r *OTPRepository) FindOTP(d *model.OTP) (*model.OTP, error) {
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id,user_id, token, created_at,uuid FROM otp WHERE token = $1"
	row := db.QueryRow(query, d.Token)

	var otp model.OTP
	err = row.Scan(&otp.Id, &otp.UserId, &otp.Token, &otp.CreatedAt, &otp.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // OTP not found, return nil without an error
		}
		return nil, err
	}

	return &otp, nil
}

func (r *OTPRepository) DeleteOTP(id int) error {
	db, err := database.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM otp WHERE id = $1"
	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
