package model

import "time"

type DailyMailCalcModel struct {
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	SubscriptionID int       `json:"subscription_id"`
	MailsForADay   int       `json:"mails_for_a_day"`
	MailsSent      int       `json:"mails_sent"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
