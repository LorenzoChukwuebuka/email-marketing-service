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
	ScheduledAt    *time.Time      `json:"scheduled_at" gorm:"index"`
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

type CampaignResponse struct {
	ID             uint                    `json:"-"`
	UUID           string                  `json:"uuid"`
	Name           string                  `json:"name"`
	Subject        *string                 `json:"subject"`
	PreviewText    *string                 `json:"preview_text"`
	SenderId       *string                 `json:"sender_id"`
	UserId         string                  `json:"user_id"`
	SenderFromName *string                 `json:"sender_from_name"`
	TemplateId     *string                 `json:"template_id"`
	SentTemplateId *string                 `json:"sent_template_id"`
	RecipientInfo  *string                 `json:"recipient_info"`
	IsPublished    bool                    `json:"is_published"`
	Status         string                  `json:"status"`
	TrackType      string                  `json:"track_type"`
	IsArchived     bool                    `json:"is_archived"`
	SentAt         *time.Time              `json:"sent_at"`
	CreatedBy      string                  `json:"created_by"`
	LastEditedBy   string                  `json:"last_edited_by"`
	Template       *string                 `json:"template"`
	Sender         *string                 `json:"sender"`
	ScheduledAt    *string                 `json:"scheduled_at"`
	CreatedAt      string                  `json:"created_at"`
	UpdatedAt      string                  `json:"updated_at"`
	DeletedAt      *string                 `json:"deleted_at"`
	CampaignGroups []CampaignGroupResponse `json:"campaign_groups"`
}

type CampaignGroupResponse struct {
	ID         uint    `json:"-"`
	UUID       string  `json:"uuid"`
	CampaignId uint    `json:"campaign_id"`
	GroupId    uint    `json:"group_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type EmailCampaignResultResponse struct {
	ID           uint       `json:"-"`
	UserID       string     `json:"user_id"`
	CampaignID   string     `json:"campaign_id"`
	Version      string     `json:"version"`
	SentAt       time.Time  `json:"sent_at"`
	OpenedAt     *time.Time `json:"opened_at,omitempty"`
	ClickedAt    *time.Time `json:"clicked_at,omitempty"`
	ConversionAt *time.Time `json:"conversion_at,omitempty"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	DeletedAt    *string    `json:"deleted_at"`
}
