package model

import "time"

type OTP struct {
	Id        int       `json:"id"`
	UUID      string    `json:"uuid"`
	UserId    int       `json:"user_id"`
	Token     string    `json:"token" validated:"required"`
	CreatedAt time.Time `json:"created_at"`
}
