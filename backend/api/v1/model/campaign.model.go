package model

import (
	"time"

	"gorm.io/gorm"
)

type CampaignStatus string

const (
	Draft     CampaignStatus = "draft"
	Saved     CampaignStatus = "saved"
	Scheduled CampaignStatus = "scheduled"
	Sent      CampaignStatus = "sent"
)

type Campaign struct {
	gorm.Model
	UUID           string          `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	Name           string          `json:"name"`
	Subject        *string         `json:"subject" gorm:"type:text"`
	PreviewText    *string         `json:"preview_text"`
	SenderId       *string         `json:"sender_id"`
	UserId         string          `json:"user_id"`
	SenderFromName *string         `json:"senderFromName"`
	TemplateId     *string         `json:"templateId"`
	SentTemplateId *string         `json:"sentTemplateId"`
	RecipientInfo  *string         `json:"recipientInfo"`
	IsPublished    bool            `json:"isPublished"`
	Status         CampaignStatus  `json:"status" gorm:"type:varchar(20);default:'draft';index"`
	TrackType      string          `json:"trackType"`
	IsArchived     bool            `json:"isArchived"`
	SentAt         *time.Time      `json:"sentAt"`
	CreatedBy      string          `json:"createdBy"`
	LastEditedBy   string          `json:"lastEditedBy"`
	Template       *string         `json:"template"`
	Sender         *string         `json:"sender"`
	CampaignGroups []CampaignGroup `json:"campaign_groups" gorm:"foreignKey:CampaignId"`
}

type CampaignGroup struct {
	gorm.Model
	UUID         string       `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	CampaignId   uint         `json:"campaign_id"`
	GroupId      uint         `json:"group_id"`
	Campaign     Campaign     `json:"-" gorm:"foreignKey:CampaignId"`
	ContactGroup ContactGroup `json:"-" gorm:"foreignKey:GroupId"`
}

type EmailCampaignResult struct {
	gorm.Model
	UserID       string `gorm:"index"`
	CampaignID   string `gorm:"index"`
	Version      string `gorm:"size:1"`
	SentAt       time.Time
	OpenedAt     *time.Time
	ClickedAt    *time.Time
	ConversionAt *time.Time
}
