package model

import (	 
	"time"
)

type DailyMailCalc struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	UUID           string    `json:"uuid"`
	SubscriptionID int       `json:"subscription_id"`
	MailsForADay   int       `json:"mails_for_a_day"`
	MailsSent      int       `json:"mails_sent"`
	RemainingMails int       `json:"remaining_mails"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type DailyMailCalcResponseModel struct {
	ID             int       `json:"-"`
	UUID           string    `json:"uuid"`
	SubscriptionID int       `json:"subscription_id"`
	MailsForADay   int       `json:"mails_for_a_day"`
	MailsSent      int       `json:"mails_sent"`
	RemainingMails int       `json:"remaining_mails"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}
