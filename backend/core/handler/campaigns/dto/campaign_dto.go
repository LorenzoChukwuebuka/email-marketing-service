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
	CompanyId  string `json:"company_id"`
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
	ID                  *string          `json:"id_ref"`
	UserID              *string          `json:"user_id"`
	CompanyID           *string          `json:"company_id"`
	Name                *string          `json:"name"`
	SenderName          *string          `json:"sender_name"`
	FromEmail           *string          `json:"from_email"`
	Subject             *string          `json:"subject"`
	Type                *string          `json:"type"`
	EmailHtml           *string          `json:"email_html"`
	EmailDesign         *json.RawMessage `json:"email_design"`
	IsEditable          bool             `json:"is_editable"`
	IsPublished         bool             `json:"is_published"`
	IsPublicTemplate    bool             `json:"is_public_template"`
	IsGalleryTemplate   bool             `json:"is_gallery_template"`
	Tags        *string          `json:"tags"`
	Description *string          `json:"description"`
	ImageUrl    *string          `json:"image_url"`
	IsActive    bool             `json:"is_active"`
	EditorType  *string          `json:"editor_type"`
	CreatedAt   *time.Time       `json:"created_at"`
	UpdatedAt   *time.Time       `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at"`
}

type EmailCampaignResultResponse struct {
	ID              string                             `json:"id"`
	CompanyID       string                             `json:"company_id"`
	CampaignID      string                             `json:"campaign_id"`
	RecipientEmail  string                             `json:"recipient_email"`
	RecipientName   *string                            `json:"recipient_name"`
	Version         *string                            `json:"version"`
	SentAt          *time.Time                         `json:"sent_at"`
	OpenedAt        *time.Time                         `json:"opened_at"`
	OpenCount       *int32                             `json:"open_count"`
	ClickedAt       *time.Time                         `json:"clicked_at"`
	ClickCount      *int32                             `json:"click_count"`
	ConversionAt    *time.Time                         `json:"conversion_at"`
	BounceStatus    *string                            `json:"bounce_status"`
	UnsubscribedAt  *time.Time                         `json:"unsubscribed_at"`
	ComplaintStatus *bool                              `json:"complaint_status"`
	DeviceType      *string                            `json:"device_type"`
	Location        *string                            `json:"location"`
	RetryCount      *int32                             `json:"retry_count"`
	Notes           *string                            `json:"notes"`
	CreatedAt       *time.Time                         `json:"created_at"`
	UpdatedAt       *time.Time                         `json:"updated_at"`
	DeletedAt       *time.Time                         `json:"deleted_at"`
	Group           []GetCampaignContactGroupsResponse `json:"group"`
}

type GetCampaignContactGroupsResponse struct {
	ID          string     `json:"id"`
	GroupName   string     `json:"group_name"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}
