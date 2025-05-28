package dto

import (
	"email-marketing-service/internal/enums"
	"time"
)

type CampaignDTO struct {
	Name           string               `json:"name" validate:"required"`
	Subject        string               `json:"subject" `
	PreviewText    string               `json:"preview_text"`
	UserId         string               `json:"user_id" validate:"required"`
	CompanyID      string               `json:"company_id" validate:"required"`
	SenderFromName string               `json:"sender_from_name"`
	TemplateId     string               `json:"template_id"`
	SentTemplateId string               `json:"sent_template_id"`
	RecipientInfo  string               `json:"recipient_info"`
	IsPublished    bool                 `json:"is_published"`
	Status         enums.CampaignStatus `json:"status" `
	TrackType      string               `json:"track_type"`
	IsArchived     bool                 `json:"is_archived"`
	SentAt         time.Time            `json:"sent_at"`
	ScheduledAt    time.Time            `json:"scheduled_at"`
	HasCustomLogo  bool                 `json:"has_custom_logo"`
	Template       string               `json:"template"`
	Sender         string               `json:"sender"`
}

type CampaignGroupDTO struct {
	CampaignId string `json:"campaign_id" validate:"required"`
	GroupId    string `json:"group_id" validate:"required"`
	UserId     string `json:"user_id"`
}

type SendCampaignDTO struct {
	CampaignId string `json:"campaign_id" validate:"required"`
	UserId     string `json:"user_id"`
}
