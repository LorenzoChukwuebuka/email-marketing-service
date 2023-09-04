package repository

import (
	"database/sql"
	"email-marketing-service/api/database"
	"email-marketing-service/api/model"
	"fmt"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(d *model.User) (*model.User, error) {

	query := "INSERT INTO users (uuid,firstname,middlename,lastname,username, email,password,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(d.UUID, d.FirstName, d.MiddleName, d.LastName, d.UserName, d.Email, d.Password, time.Now()).Scan(&d.ID); err != nil {
		return nil, err
	}

	return d, nil
}

func (r *UserRepository) CheckIfEmailAlreadyExists(d *model.User) (bool, error) {

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	var exists bool
	err := r.DB.QueryRow(query, d.Email).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) VerifyUserAccount(d *model.User) error {
	query := "UPDATE users SET verified = $2, verified_at = $3 WHERE id = $1"
	_, err := r.DB.Exec(query, d.ID, d.Verified, d.VerifiedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Login(d *model.User) (*model.UserResponse, error) {

	// query := "SELECT * FROM users WHERE email = $1 AND verified = true"
	query := "SELECT id, uuid, firstname, middlename, lastname, username, email, password, verified, verified_at FROM users WHERE email = $1 AND verified = true"
	row := r.DB.QueryRow(query, d.Email)

	var user model.UserResponse

	var verifiedAt sql.NullTime

	err := row.Scan(&user.ID, &user.UUID, &user.FirstName, &user.MiddleName, &user.LastName, &user.UserName, &user.Email, &user.Password, &user.Verified, &verifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found: %w", err) // User not found, return nil without an error
		}
		return nil, err
	}

	if verifiedAt.Valid {
		user.VerifiedAt = verifiedAt.Time.Format(time.RFC3339Nano)
	}

	return &user, nil
}

func (r *UserRepository) FindUserById(d *model.User) (*model.User, error) {

	query := "SELECT id, uuid, firstname, middlename, lastname, username, email, password, verified, created_at, verified_at, updated_at, deleted_at FROM users WHERE id = $1"
	row := r.DB.QueryRow(query, d.ID)

	err := row.Scan(&d.ID,
		&d.UUID,
		&d.FirstName,
		&d.MiddleName,
		&d.LastName,
		&d.UserName,
		&d.Email,
		&d.Password,
		&d.Verified,
		&d.CreatedAt,
		&d.VerifiedAt,
		&d.UpdatedAt,
		&d.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // User not found, return nil without an error
		}
		return nil, err
	}

	return d, nil
}

func (r *UserRepository) FindUserByEmail(d *model.User) (*model.User, error) {

	query := "SELECT id, uuid, firstname, middlename, lastname, username, email, password, verified, created_at, verified_at, updated_at, deleted_at FROM users WHERE id = $1"
	row := r.DB.QueryRow(query, d.Email)

	err := row.Scan(&d.ID, &d.UserName, &d.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // User not found, return nil without an error
		}
		return nil, err
	}

	return d, nil
}

func (r *UserRepository) ResetPassword(d *model.User) error {
	db, err := database.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE users SET password = $1 WHERE id = $2"

	_, err = db.Exec(query, d.Password, d.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindAllUsers() ([]model.UserResponse, error) {

	query := "SELECT id, uuid, firstname, middlename, lastname, username, email, password, verified, created_at, verified_at, updated_at, deleted_at FROM users"

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserResponse

	for rows.Next() {
		var user model.UserResponse
		var verifiedAt, updatedAt sql.NullTime
		err := rows.Scan(
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
			&verifiedAt,
			&updatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if verifiedAt.Valid {
			user.VerifiedAt = verifiedAt.Time.Format(time.RFC3339Nano)
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time.Format(time.RFC3339Nano)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) ChangeUserPassword(d *model.User) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(d.Password, d.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUserRecords(d *model.User) error {
	query := "UPDATE users  SET firstname = $1,middlename = $2,lastname = $3, username = $4,updated_at=$5 WHERE id = $6"

	stmt, err := r.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(d.FirstName, d.MiddleName, d.LastName, d.UserName, time.Now(), d.ID)

	if err != nil {
		return err
	}

	return nil
}
