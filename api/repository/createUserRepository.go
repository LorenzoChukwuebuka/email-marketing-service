package repository

import (
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
	query := "INSERT INTO users (uuid,firstname,middlename,lastname,username, email,password,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"
	_, err = db.Exec(query, d.UUID, d.FirstName, d.MiddleName, d.LastName, d.UserName, d.Email, d.Password, d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func CheckIfEmailAlreadyExists(d *model.User) {

}
