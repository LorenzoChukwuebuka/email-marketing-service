package model

import (
	"database/sql"
	"time"
)

type UserSessionModelStruct struct {
	Id        int          `json:"id"`
	UUID      string       `json:"uuid"`
	UserId    int          `json:"user_id"`
	Device    *string      `json:"device"`
	IPAddress *string      `json:"ip_address"`
	Location  *string      `json:"location"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
