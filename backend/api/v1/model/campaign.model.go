package model

import (
	"gorm.io/gorm"
	"time"
)

type CampaignStatus string

type CampaignTrackType string

const (
	Draft       CampaignStatus = "draft"
	Saved       CampaignStatus = "saved"
	Scheduled   CampaignStatus = "scheduled"
	Sent        CampaignStatus = "sent"
	FailedC     CampaignStatus = "failed"
	Proccessing CampaignStatus = "proccessing"
)

const (
	Track CampaignTrackType = "track"
)

type Campaign struct {
	gorm.Model
	UUID           string            `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	Name           string            `json:"name"`
	Subject        *string           `json:"subject" gorm:"type:text"`
	PreviewText    *string           `json:"preview_text"`
	UserId         string            `json:"user_id"`
	SenderFromName *string           `json:"sender_from_name"`
	TemplateId     *uint             `json:"template_id"`
	Template       *Template         `json:"template"`
	SentTemplateId *string           `json:"sent_template_id"`
	RecipientInfo  *string           `json:"recipient_info"`
	IsPublished    bool              `json:"is_published"`
	Status         CampaignStatus    `json:"status" gorm:"type:varchar(20);default:'draft';index"`
	TrackType      CampaignTrackType `json:"trackType"`
	IsArchived     bool              `json:"isArchived"`
	SentAt         *time.Time        `json:"sentAt"`
	Sender         *string           `json:"sender"`
	ScheduledAt    *time.Time        `json:"scheduled_at" gorm:"index"`
	HasCustomLogo  bool              `json:"has_custom_logo"`
	CampaignGroups []CampaignGroup   `json:"campaign_groups" gorm:"foreignKey:CampaignId"`
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
	CampaignID      string     `json:"campaign_id" gorm:"index"`                     // Campaign identifier
	RecipientEmail  string     `json:"recipient_email" gorm:"size:255;index"`        // Email of the recipient
	RecipientName   *string    `json:"recipient_name" gorm:"size:255;default:null"`  // Name of the recipient
	Version         string     `json:"version" gorm:"size:10"`                       // Version of the campaign (e.g., A/B testing)
	SentAt          time.Time  `json:"sent_at" gorm:"type:TIMESTAMP;default:null"`   // Timestamp when the email was sent
	OpenedAt        *time.Time `json:"opened_at" gorm:"type:TIMESTAMP;default:null"` // Timestamp when the email was opened
	OpenCount       int        `json:"open_count" gorm:"default:0"`
	ClickedAt       *time.Time `json:"clicked_at" gorm:"type:TIMESTAMP;default:null"` // Timestamp when a link in the email was clicked
	ClickCount      int        `gorm:"default:0"`
	ConversionAt    *time.Time `json:"conversion_at" gorm:"type:TIMESTAMP;default:null"`   // Timestamp when a conversion occurred (e.g., purchase)
	BounceStatus    string     `gorm:"size:20"`                                            // Status if the email bounced (e.g., "soft", "hard")
	UnsubscribeAt   *time.Time `json:"unsubscribed_at" gorm:"type:TIMESTAMP;default:null"` // Timestamp when the recipient unsubscribed
	ComplaintStatus bool       `json:"complaint_status"`                                   // Whether the recipient marked the email as spam
	DeviceType      string     `gorm:"size:50"`                                            // Type of device used to open the email
	Location        string     `gorm:"size:100"`                                           // Geographic location of the recipient (if tracked)
	RetryCount      int        `json:"retry_count"`                                        // Number of retry attempts for sending the email
	Notes           string     `gorm:"type:text"`                                          // Additional notes or metadata
}

type CampaignResponse struct {
	ID             uint                    `json:"-"`
	UUID           string                  `json:"uuid"`
	Name           string                  `json:"name"`
	Subject        *string                 `json:"subject"`
	PreviewText    *string                 `json:"preview_text"`
	UserId         string                  `json:"user_id"`
	SenderFromName *string                 `json:"sender_from_name"`
	TemplateId     *uint                   `json:"template_id"`
	SentTemplateId *string                 `json:"sent_template_id"`
	RecipientInfo  *string                 `json:"recipient_info"`
	IsPublished    bool                    `json:"is_published"`
	Status         string                  `json:"status"`
	TrackType      string                  `json:"track_type"`
	IsArchived     bool                    `json:"is_archived"`
	SentAt         *time.Time              `json:"sent_at"`
	CreatedBy      string                  `json:"created_by"`
	LastEditedBy   string                  `json:"last_edited_by"`
	Template       *Template               `json:"template"`
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
	CampaignID   string     `json:"campaign_id"`
	Version      string     `json:"version"`
	SentAt       time.Time  `json:"sent_at"`
	OpenedAt     *time.Time `json:"opened_at"`
	OpenCount    int        `json:"open_count"`
	ClickedAt    *time.Time `json:"clicked_at"`
	ClickCount   int        `json:"click_count"`
	ConversionAt *time.Time `json:"conversion_at"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	DeletedAt    *string    `json:"deleted_at"`
}
