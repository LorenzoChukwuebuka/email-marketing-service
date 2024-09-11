package model

import "gorm.io/gorm"

type EmailBox struct {
	gorm.Model
	UUID    string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserName  string `json:"user_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content []byte `json:"content"`
	Mailbox string `json:"mailbox"`
}
