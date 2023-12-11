package model

import (
	"database/sql"
	"time"
)

type DailyMailCalcModel struct {
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	SubscriptionID int       `json:"subscription_id"`
	MailsForADay   int       `json:"mails_for_a_day"`
	MailsSent      int       `json:"mails_sent"`
	RemainingMails int       `json:"remaining_mails"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt     sql.NullTime `json:"updated_at"`
}


type DailyMailCalcResponseModel struct {
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	SubscriptionID int       `json:"subscription_id"`
	MailsForADay   int       `json:"mails_for_a_day"`
	MailsSent      int       `json:"mails_sent"`
	RemainingMails int       `json:"remaining_mails"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
