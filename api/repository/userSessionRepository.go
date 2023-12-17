package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"fmt"
)

type UserSessionRepository struct {
	DB *sql.DB
}

func NewUserSessionRepository(db *sql.DB) *UserSessionRepository {
	return &UserSessionRepository{
		DB: db,
	}
}

func (r *UserSessionRepository) CreateSession(session *model.UserSessionModelStruct) error {
	query := `
		INSERT INTO user_sessions (uuid, user_id, device, ip_address, browser, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	var sessionID int
	err = stmt.QueryRow(
		session.UUID,
		session.UserId,
		session.Device,
		session.IPAddress,
		session.Browser,
		session.CreatedAt,
	).Scan(&sessionID)

	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}

	return nil
}

func (r *UserSessionRepository) GetSessionsByUserID(userID int) ([]model.UserSessionResponseModel, error) {
	query := `
		SELECT id, uuid, user_id, device, ip_address, browser, created_at, updated_at, deleted_at
		FROM user_sessions
		WHERE user_id = $1
	`

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing select statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("error executing select query: %v", err)
	}
	defer rows.Close()

	var sessions []model.UserSessionResponseModel
	for rows.Next() {
		var session model.UserSessionResponseModel

		var updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&session.Id,
			&session.UUID,
			&session.UserId,
			&session.Device,
			&session.IPAddress,
			&session.Browser,
			&session.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		SetTime(updatedAt, &session.UpdatedAt)
		SetTime(deletedAt, &session.DeletedAt)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *UserSessionRepository) DeleteSession(sessionId string) error {

	query := "DELETE FROM user_sessions WHERE uuid = $1"
	_, err := r.DB.Exec(query, sessionId)
	if err != nil {
		return err
	}
	return nil
	 
}
