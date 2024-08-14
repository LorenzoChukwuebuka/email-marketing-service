package dto

import (
	"encoding/json"
)

type TemplateDTO struct {
	UserId            string           `json:"user_id" validate:"required"`
	TemplateName      string           `json:"templateName" validate:"required"`
	SenderName        *string          `json:"senderName"`
	FromEmail         *string          `json:"fromEmail"`
	Subject           *string          `json:"subject"`
	Type              string           `json:"type" validate:"required"`
	EmailHtml         string           `json:"emailHtml"`
	EmailDesign       *json.RawMessage `json:"emailDesign"`
	IsEditable        bool             `json:"isEditable"`
	IsPublished       bool             `json:"isPublished"`
	IsPublicTemplate  bool             `json:"isPublicTemplate"`
	IsGalleryTemplate bool             `json:"isGalleryTemplate"`
	Tags              string           `json:"tags"`
	Description       *string          `json:"description"`
	ImageUrl          *string          `json:"imageUrl"`
	IsActive          bool             `json:"isActive"`
	EditorType        *string          `json:"editorType"`
}
