package model

import (
	"gorm.io/gorm"
)

type DailyMailCalc struct {
	gorm.Model
	UUID           string `json:"uuid"`
	SubscriptionID int    `json:"subscription_id"`
	MailsForADay   int    `json:"mails_for_a_day"`
	MailsSent      int    `json:"mails_sent"`
	RemainingMails int    `json:"remaining_mails"`
}

type DailyMailCalcResponseModel struct {
	ID             int    `json:"-"`
	UUID           string `json:"uuid"`
	SubscriptionID int    `json:"subscription_id"`
	MailsForADay   int    `json:"mails_for_a_day"`
	MailsSent      int    `json:"mails_sent"`
	RemainingMails int    `json:"remaining_mails"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
