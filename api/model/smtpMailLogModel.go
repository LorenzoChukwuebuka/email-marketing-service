package model

import "time"

type EmailStatus string

const (
	Sending EmailStatus = "sending"
	Failed  EmailStatus = "failed"
	Success EmailStatus = "success"
)

type SentEmails struct {
	Id             uint        `json:"-" gorm:"primaryKey;index"`
	UUID           string      `json:"uuid" gorm:"type:uuid;index"`
	Sender         uint        `json:"sender_id" gorm:"index"`
	Recipient      uint        `json:"recipient"`
	MessageContent string      `json:"message_content"`
	Status         EmailStatus `json:"status"`
	CreatedAt      time.Time   `json:"created_at" gorm:"type:TIMESTAMP;null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      *time.Time  `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
}
