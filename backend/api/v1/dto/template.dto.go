package dto

import (
	"encoding/json"
)

type TemplateDTO struct {
	UserId            string           `json:"user_id" validate:"required"`
	TemplateName      string           `json:"template_name" validate:"required"`
	SenderName        *string          `json:"sender_name"`
	FromEmail         *string          `json:"from_email"`
	Subject           *string          `json:"subject"`
	Type              string           `json:"type" validate:"required"`
	EmailHtml         string           `json:"email_html"`
	EmailDesign       *json.RawMessage `json:"email_design"`
	IsEditable        bool             `json:"is_editable"`
	IsPublished       bool             `json:"is_published"`
	IsPublicTemplate  bool             `json:"is_public_template"`
	IsGalleryTemplate bool             `json:"is_gallery_template"`
	Tags              string           `json:"tags"`
	Description       *string          `json:"description"`
	ImageUrl          *string          `json:"image_Url"`
	IsActive          bool             `json:"is_active"`
	EditorType        *string          `json:"editor_type"`
}

type SendTestMailDTO struct {
	UserId       string `json:"user_id" validate:"required"`
	EmailAddress string `json:"email_address" validate:"required"`
	TemplateId   string `json:"template_id" validate:"required"`
	Subject      string `json:"subject" validate:"required"`
}


