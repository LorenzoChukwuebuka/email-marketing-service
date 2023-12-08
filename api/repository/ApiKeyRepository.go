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

func (r *APIKeyRepository) GetUserAPIKeyByUserId(userId int) (*model.APIKeyResponseModel, error) {
	query := `SELECT id, uuid, user_id, api_key, created_at FROM api_keys WHERE user_id=$1;`

	row := r.DB.QueryRow(query, userId)

	var apiKeyModel model.APIKeyResponseModel
	err := row.Scan(&apiKeyModel.Id, &apiKeyModel.UUID, &apiKeyModel.UserId, &apiKeyModel.APIKey, &apiKeyModel.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found, return nil without an error
		}
		return nil, err
	}

	return &apiKeyModel, nil
}

func (r *APIKeyRepository) UpdateAPIKey(d *model.APIKeyModel) error {
	query := `UPDATE api_keys SET api_key = $1 WHERE user_id = $2;`

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//d.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	_, err = stmt.Exec(d.APIKey, d.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *APIKeyRepository) CheckIfAPIKEYExists(apiKey string) (bool, error) {

	query := "SELECT EXISTS(SELECT 1 FROM api_keys WHERE api_key = $1)"

	var exists bool
	err := r.DB.QueryRow(query, apiKey).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}
