package repository

import (
	"database/sql"
	"email-marketing-service/api/database"
	"email-marketing-service/api/model"
)

func CreateOTP(d *model.OTP) error {

	db, err := database.InitDB()
	if err != nil {
		return nil
	}
	defer db.Close()

	query := "Insert into otp (user_id,token,created_at,uuid)Values($1,$2,$3,$4)"

	_, err = db.Exec(query, d.UserId, d.Token, d.CreatedAt, d.UUID)

	if err != nil {
		return err
	}

	return nil

}

func FindOTP(d *model.OTP) (*model.OTP, error) {
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT user_id, token, created_at,uuid FROM otp WHERE token = $1"
	row := db.QueryRow(query, d.Token)

	var otp model.OTP
	err = row.Scan(&otp.UserId, &otp.Token, &otp.CreatedAt, &otp.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // OTP not found, return nil without an error
		}
		return nil, err
	}

	return &otp, nil
}
