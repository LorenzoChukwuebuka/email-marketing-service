package model

import "time"

type OTP struct {
	UUID     string    `json:"uuid"`
	UserId   int       `json:"user_id"`
	Token    string    `json:"token"`
	CreateAt time.Time `json:"created_at"`
}
