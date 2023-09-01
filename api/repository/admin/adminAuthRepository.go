package adminrepository

import (
	"database/sql"
	adminmodel "email-marketing-service/api/model/admin"
	"fmt"
)

type AdminRepository struct {
	DB *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) Login(d *adminmodel.AdminLogin) (*adminmodel.AdminResponse, error) {
	// query := "SELECT * FROM users WHERE email = $1 AND verified = true"
	query := "SELECT id, firstname, middlename, lastname,  email, password, type FROM admin WHERE email = $1 "
	row := r.DB.QueryRow(query, d.Email)

	var admin adminmodel.AdminResponse
	err := row.Scan(&admin.ID, &admin.FirstName, &admin.MiddleName, &admin.LastName, &admin.Email, &admin.Password, &admin.Type)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found: %w", err) // User not found, return nil without an error
		}
		return nil, err
	}

	return &admin, nil

}

func (r *AdminRepository) ChangePassword() {

}
