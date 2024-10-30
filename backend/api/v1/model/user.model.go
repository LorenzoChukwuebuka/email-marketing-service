package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UUID                 string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FullName             string     `json:"fullname" `
	Company              string     `json:"company" `
	Email                string     `json:"email"  gorm:"index"`
	PhoneNumber          string     `json:"phonenumber" gorm:"type:varchar(255);default:null"`
	Password             string     `json:"password"  gorm:"index"`
	Verified             bool       `json:"verified"`
	Blocked              bool       `json:"blocked" gorm:"default:false"`
	VerifiedAt           *time.Time `json:"verified_at" gorm:"type:TIMESTAMP;null;default:null"`
	Status               string     `json:"status" gorm:"type:varchar(20);default:active"`
	ScheduledForDeletion bool       `gorm:"default:false" json:"scheduled_for_deletion"`
	ScheduledDeletionAt  *time.Time `json:"scheduled_deletion_at"`
	LastLoginAt          *time.Time `json:"last_login_at"`
}

type UserResponse struct {
	ID                   uint    `json:"-"`
	UUID                 string  `json:"uuid"`
	FullName             string  `json:"fullname"`
	Email                string  `json:"email"`
	Company              string  `json:"company"`
	PhoneNumber          string  `json:"phonenumber"`
	Password             string  `json:"-"`
	Verified             bool    `json:"verified"`
	Blocked              bool    `json:"blocked"`
	Status               bool    `json:"status"`
	CreatedAt            string  `json:"created_at"`
	VerifiedAt           *string `json:"verified_at"`
	ScheduledForDeletion bool    `json:"scheduled_for_deletion"`
	ScheduledDeletionAt  *string `json:"scheduled_deletion_at"`
	LastLoginAt          *string `json:"last_login_at"`
	UpdatedAt            string  `json:"updated_at"`
	DeletedAt            *string `json:"deleted_at"`
}

// UserArchive model for keeping records of deleted users
type UserArchive struct {
	gorm.Model
	UserID           string       `gorm:"type:uuid" json:"user_id"`
	Email            string       `json:"email"`
	FullName         string       `json:"full_name"`
	Company          string       `json:"company"`
	DeletedAt        time.Time    `json:"deleted_at"`
	AccountCreatedAt time.Time    `json:"account_created_at"`
	VerifiedAt       *time.Time   `json:"verified_at"`
	LastLoginAt      *time.Time   `json:"last_login_at"`
	DeletionReason   string       `json:"deletion_reason"`
	AccountStats     AccountStats `gorm:"embedded"`
}

type AccountStats struct {
	TotalContacts      int64 `json:"total_contacts"`
	TotalCampaigns     int64 `json:"total_campaigns"`
	TotalTemplates     int64 `json:"total_templates"`
	TotalCampaignsSent int64 `json:"total_campaigns_sent"`
	TotalGroups        int64 `json:"total_groups"`
}

type UserTempEmail struct {
	gorm.Model
	TemporaryEmail string `gorm:"unique;not null"`
	UserId         string
}
