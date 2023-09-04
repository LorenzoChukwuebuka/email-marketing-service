package adminrepository

import (
	"database/sql"
	adminmodel "email-marketing-service/api/model/admin"
	"fmt"
	"time"
)

type AdminRepository struct {
	DB *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) CreateAdmin(d *adminmodel.AdminModel) (*adminmodel.AdminModel, error) {
	query := `
	INSERT INTO admin(firstname, middlename, lastname, email, password, type, created_at)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow( d.FirstName, d.MiddleName, d.LastName, d.Email, d.Password, "admin", time.Now()).Scan(&d.ID); err != nil {
		return nil, err
	}

	return d, nil
}

func (r *AdminRepository) Login(d *adminmodel.AdminLogin) (*adminmodel.AdminResponse, error) {

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
