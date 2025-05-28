package dto

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type TemplateResponse struct {
	ID                uuid.UUID               `json:"id"`
	UserID            uuid.UUID               `json:"user_id"`
	CompanyID         uuid.UUID               `json:"company_id"`
	TemplateName      string                  `json:"template_name"`
	SenderName        string                  `json:"sender_name"`
	FromEmail         string                  `json:"from_email"`
	Subject           string                  `json:"subject"`
	Type              string                  `json:"type"`
	EmailHtml         string                  `json:"email_html"`
	EmailDesign       json.RawMessage         `json:"email_design"`
	IsEditable        bool                    `json:"is_editable"`
	IsPublished       bool                    `json:"is_published"`
	IsPublicTemplate  bool                    `json:"is_public_template"`
	IsGalleryTemplate bool                    `json:"is_gallery_template"`
	Tags              string                  `json:"tags"`
	Description       string                  `json:"description"`
	ImageUrl          string                  `json:"image_url"`
	IsActive          bool                    `json:"is_active"`
	EditorType        string                  `json:"editor_type"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
	DeletedAt         time.Time               `json:"deleted_at"`
	User              TemplateUserResponse    `json:"user"`
	Company           TemplateCompanyResponse `json:"company"`
}

type TemplateUserResponse struct {
	UserFullname string `json:"user_fullname"`
	UserEmail    string `json:"user_email"`
	UserPicture  string `json:"user_picture"`
}

type TemplateCompanyResponse struct {
	CompanyName string `json:"company_name"`
}
