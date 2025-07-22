package dto

import (
	"encoding/json"
)

type AdminTemplateDTO struct {
	UserId            string          `json:"user_id"`
	CompanyID         string          `json:"company_id"`
	TemplateName      string          `json:"template_name" validate:"required"`
	SenderName        string          `json:"sender_name"`
	FromEmail         string          `json:"from_email"`
	Subject           string          `json:"subject"`
	Type              string          `json:"type" validate:"required"`
	EmailHtml         string          `json:"email_html"`
	EmailDesign       json.RawMessage `json:"email_design"`
	IsEditable        bool            `json:"is_editable"`
	IsPublished       bool            `json:"is_published"`
	IsPublicTemplate  bool            `json:"is_public_template"`
	IsGalleryTemplate bool            `json:"is_gallery_template"`
	Tags              string          `json:"tags"`
	Description       string          `json:"description"`
	ImageUrl          string          `json:"image_Url"`
	IsActive          bool            `json:"is_active"`
	EditorType        string          `json:"editor_type"`
	TemplateID        string          `json:"template_id"`
}


type AdminFetchGalleryTemplatesDTO struct {
	TemplateId string `json:"template_id"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Search     string `json:"search"`
}