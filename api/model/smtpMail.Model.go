package model

import "time"

type EmailStatus string

const (
	Sending EmailStatus = "sending"
	Failed  EmailStatus = "failed"
	Success EmailStatus = "success"
)

type SentEmails struct {
	Id             uint
	UUID           string
	Sender         uint
	Recipient      uint
	MessageContent string
	Status         EmailStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
