package model

import (
    "gorm.io/gorm"
    "time"
)

type MailUsage struct {
    gorm.Model
    UUID           string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
    SubscriptionID int       `json:"subscription_id"`
    PeriodStart    time.Time `json:"period_start"`
    PeriodEnd      time.Time `json:"period_end"`
    LimitAmount    int       `json:"limit_amount"`
    MailsSent      int       `json:"mails_sent"`
    RemainingMails int       `json:"remaining_mails"`
}

type MailUsageResponseModel struct {
    ID             int       `json:"-"`
    UUID           string    `json:"uuid"`
    SubscriptionID int       `json:"subscription_id"`
    PeriodStart    time.Time `json:"period_start"`
    PeriodEnd      time.Time `json:"period_end"`
    LimitAmount    int       `json:"limit_amount"`
    MailsSent      int       `json:"mails_sent"`
    RemainingMails int       `json:"remaining_mails"`
    CreatedAt      string    `json:"created_at"`
    UpdatedAt      string    `json:"updated_at"`
}