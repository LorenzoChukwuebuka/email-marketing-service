package model

import (
	"database/sql"
	"time"
)

type APIKeyModel struct {
	Id      int          `json:"id"`
	UUID      string       `json:"uuid"`
	UserId    int          `json:"user_id"`
	APIKey    string       `json:"api_key"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}


type APIKeyResponseModel struct{
	Id      int          `json:"id"`
	UUID      string       `json:"uuid"`
	UserId    int          `json:"user_id"`
	APIKey    string       `json:"api_key"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
