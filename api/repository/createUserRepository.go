package repository

import (
	"database/sql"
	"email-marketing-service/api/database"
	"email-marketing-service/api/model"
	"fmt"
)

func CreateUser(d *model.User) (*model.User, error) {
	// Initialize the database connection
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := "INSERT INTO users (uuid,firstname,middlename,lastname,username, email,password) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	err = db.QueryRow(query, d.UUID, d.FirstName, d.MiddleName, d.LastName, d.UserName, d.Email, d.Password).Scan(&d.ID)

	if err != nil {
		return nil, err
	}
	return d, nil
}

func CheckIfEmailAlreadyExists(d *model.User) (bool, error) {
	var user model.User
	db, err := database.InitDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := "SELECT * FROM users WHERE email = $1"
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	err = db.QueryRow(query, d.Email).Scan(
		&user.ID,
		&user.UUID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.VerifiedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return false, nil
		}
		fmt.Println("Error:", err)
		return false, err
	}

	// Email exists, return true
	return true, nil
}
