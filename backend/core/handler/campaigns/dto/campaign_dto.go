package dto

import (
	"email-marketing-service/internal/enums"
	"encoding/json"
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
	CompanyId string `json:"company_id"`
}

type FetchCampaignDTO struct {
	UserID      string `json:"user_id"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	SearchQuery string `json:"search_query"`
	CompanyID   string `json:"company_id"`
	CampaignID  string `json:"campaign_id"`
}

type CampaignResponseDTO struct {
	ID             string            `json:"id"`
	CompanyID      string            `json:"company_id"`
	Name           string            `json:"name"`
	Subject        *string           `json:"subject"`
	PreviewText    *string           `json:"preview_text"`
	UserID         string            `json:"user_id"`
	SenderFromName *string           `json:"sender_from_name"`
	TemplateID     *string           `json:"template_id"`
	SentTemplateID *string           `json:"sent_template_id"`
	RecipientInfo  *string           `json:"recipient_info"`
	IsPublished    bool              `json:"is_published"`
	Status         *string           `json:"status"`
	TrackType      *string           `json:"track_type"`
	IsArchived     bool              `json:"is_archived"`
	SentAt         *time.Time        `json:"sent_at"`
	Sender         *string           `json:"sender"`
	ScheduledAt    *time.Time        `json:"scheduled_at"`
	HasCustomLogo  bool              `json:"has_custom_logo"`
	CreatedAt      *time.Time        `json:"created_at"`
	UpdatedAt      *time.Time        `json:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at"`
	User           UserResponse      `json:"user"`
	Company        CompanyResponse   `json:"company"`
	Template       *TemplateResponse `json:"template"`
}

type UserResponse struct {
	UserID_2        string     `json:"user_id_2"`
	UserFullname    string     `json:"user_fullname"`
	UserEmail       string     `json:"user_email"`
	UserPhonenumber *string    `json:"user_phonenumber"`
	UserPicture     *string    `json:"user_picture"`
	UserVerified    bool       `json:"user_verified"`
	UserBlocked     bool       `json:"user_blocked"`
	UserVerifiedAt  *time.Time `json:"user_verified_at"`
	UserStatus      string     `json:"user_status"`
	UserLastLoginAt *time.Time `json:"user_last_login_at"`
	UserCreatedAt   time.Time  `json:"user_created_at"`
	UserUpdatedAt   time.Time  `json:"user_updated_at"`
}

type CompanyResponse struct {
	CompanyIDRef     string    `json:"company_id_ref"`
	CompanyName      *string   `json:"company_name"`
	CompanyCreatedAt time.Time `json:"company_created_at"`
	CompanyUpdatedAt time.Time `json:"company_updated_at"`
}

type TemplateResponse struct {
	TemplateIDRef             *string          `json:"template_id_ref"`
	TemplateUserID            *string          `json:"template_user_id"`
	TemplateCompanyID         *string          `json:"template_company_id"`
	TemplateName              *string          `json:"template_name"`
	TemplateSenderName        *string          `json:"template_sender_name"`
	TemplateFromEmail         *string          `json:"template_from_email"`
	TemplateSubject           *string          `json:"template_subject"`
	TemplateType              *string          `json:"template_type"`
	TemplateEmailHtml         *string          `json:"template_email_html"`
	TemplateEmailDesign       *json.RawMessage `json:"template_email_design"`
	TemplateIsEditable        bool             `json:"template_is_editable"`
	TemplateIsPublished       bool             `json:"template_is_published"`
	TemplateIsPublicTemplate  bool             `json:"template_is_public_template"`
	TemplateIsGalleryTemplate bool             `json:"template_is_gallery_template"`
	TemplateTags              *string          `json:"template_tags"`
	TemplateDescription       *string          `json:"template_description"`
	TemplateImageUrl          *string          `json:"template_image_url"`
	TemplateIsActive          bool             `json:"template_is_active"`
	TemplateEditorType        *string          `json:"template_editor_type"`
	TemplateCreatedAt         *time.Time       `json:"template_created_at"`
	TemplateUpdatedAt         *time.Time       `json:"template_updated_at"`
	TemplateDeletedAt         *time.Time       `json:"template_deleted_at"`
}
