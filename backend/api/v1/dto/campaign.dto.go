package dto

import (
	"time"
)

type CampaignStatus string

const (
	Draft     CampaignStatus = "draft"
	Saved     CampaignStatus = "saved"
	Scheduled CampaignStatus = "scheduled"
	Sent      CampaignStatus = "sent"
)

type CampaignDTO struct {
	Name           string         `json:"name" validate:"required"`
	Subject        *string        `json:"subject" `
	PreviewText    *string        `json:"preview_text"`
	SenderId       *string        `json:"sender_id"`
	UserId         string         `json:"user_id" validate:"required"`
	SenderFromName *string        `json:"senderFromName"`
	TemplateId     *string        `json:"templateId"`
	SentTemplateId *string        `json:"sentTemplateId"`
	RecipientInfo  *string        `json:"recipientInfo"`
	IsPublished    bool           `json:"isPublished"`
	Status         CampaignStatus `json:"status" `
	TrackType      string         `json:"trackType"`
	IsArchived     bool           `json:"isArchived"`
	SentAt         *time.Time     `json:"sentAt"`
	CreatedBy      string         `json:"createdBy"`
	LastEditedBy   string         `json:"lastEditedBy"`
	Template       *string        `json:"template"`
	Sender         *string        `json:"sender"`
}
