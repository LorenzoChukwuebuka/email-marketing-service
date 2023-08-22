package repository

import (
	"database/sql"
	"email-marketing-service/api/database"
	"email-marketing-service/api/model"
)

func CreateUser(d *model.User) (*model.User, error) {
	// Initialize the database connection
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := "INSERT INTO users (uuid,firstname,middlename,lastname,username, email,password) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
	err = db.QueryRow(query, d.UUID, d.FirstName, d.MiddleName, d.LastName, d.UserName, d.Email, d.Password).Scan(&d.ID)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return d, nil
}

func CheckIfEmailAlreadyExists(d *model.User) (bool, error) {
	db, err := database.InitDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	var exists bool
	err = db.QueryRow(query, d.Email).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

func FindUserById(d *model.User) (*model.User, error) {
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, username, email FROM users WHERE id = $1"
	row := db.QueryRow(query, d.ID)

	err = row.Scan(&d.ID, &d.UserName, &d.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // User not found, return nil without an error
		}
		return nil, err
	}

	return d, nil
}
