package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type TemplateType string

const (
	Transactional TemplateType = "Transactional"
	Marketing     TemplateType = "Marketing"
)

type Template struct {
	gorm.Model
	UUID              string           `gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId            string           `json:"user_id"`
	TemplateName      string           `json:"template_name"`
	SenderName        *string          `json:"sender_name" gorm:"default:null"`
	FromEmail         *string          `json:"from_email" gorm:"default:null"`
	Subject           *string          `json:"subject" gorm:"default:null"`
	Type              TemplateType     `json:"type"`
	EmailHtml         string           `json:"email_html"`
	EmailDesign       *json.RawMessage `json:"email_design" gorm:"type:jsonb;default:null"`
	IsEditable        bool             `json:"is_editable"`
	IsPublished       bool             `json:"is_published"`
	IsPublicTemplate  bool             `json:"is_public_template"`
	IsGalleryTemplate bool             `json:"is_gallery_template"`
	Tags              string           `json:"tags"`
	Description       *string          `json:"description" gorm:"default:null"`
	ImageUrl          *string          `json:"image_Url" gorm:"default:null"`
	IsActive          bool             `json:"is_active"`
	EditorType        *string          `json:"editor_type" gorm:"default:null"`
}

type TemplateResponse struct {
	ID                uint             `json:"id"`
	UUID              string           `json:"uuid"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
	DeletedAt         *string          `json:"deleted_at"`
	UserId            string           `json:"user_id"`
	TemplateName      string           `json:"template_name"`
	SenderName        *string          `json:"sender_name"`
	FromEmail         *string          `json:"from_email"`
	Subject           *string          `json:"subject"`
	Type              TemplateType     `json:"type"`
	EmailHtml         string           `json:"email_html"`
	EmailDesign       *json.RawMessage `json:"email_design"`
	IsEditable        bool             `json:"is_editable"`
	IsPublished       bool             `json:"is_published"`
	IsPublicTemplate  bool             `json:"is_public_template"`
	IsGalleryTemplate bool             `json:"is_gallery_template"`
	Tags              string           `json:"tags"`
	Description       *string          `json:"description"`
	ImageUrl          *string          `json:"image_url"`
	IsActive          bool             `json:"is_active"`
	EditorType        *string          `json:"editor_type"`
}
