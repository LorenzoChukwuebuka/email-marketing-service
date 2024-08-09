package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Email struct {
	gorm.Model
	TeamID            string          `json:"teamId"`
	TeamMemberID      string          `json:"teamMemberId"`
	TemplateID        int             `json:"templateId"`
	TemplateName      string          `json:"templateName"`
	CampaignID        *int            `json:"campaignId"`
	SenderName        *string         `json:"senderName"`
	FromEmail         *string         `json:"fromEmail"`
	Subject           *string         `json:"subject"`
	Type              string          `json:"type"`
	EmailHtml         string          `json:"emailHtml" gorm:"type:jsonb"`
	EmailDesign       json.RawMessage `json:"emailDesign"`
	IsEditable        bool            `json:"isEditable"`
	IsPublished       bool            `json:"isPublished"`
	IsPublicTemplate  bool            `json:"isPublicTemplate"`
	IsGalleryTemplate bool            `json:"isGalleryTemplate"`
	Tags              string          `json:"tags"`
	Description       *string         `json:"description"`
	ImageUrl          *string         `json:"imageUrl"`
	IsActive          bool            `json:"isActive"`
	EditorType        *string         `json:"editorType"`
}
