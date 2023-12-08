package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
)

type APIKeyRepository struct {
	DB *sql.DB
}

func NewAPIkeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{
		DB: db,
	}
}

func (r *APIKeyRepository) CreateAPIKey(d *model.APIKeyModel) (*model.APIKeyModel, error) {
	query := `INSERT INTO api_keys(uuid, user_id, api_key, created_at) VALUES ($1, $2, $3, $4) RETURNING id;`

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var insertedID int
	err = stmt.QueryRow(d.UUID, d.UserId, d.APIKey, d.CreatedAt).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	d.Id = insertedID
	return d, nil

}
