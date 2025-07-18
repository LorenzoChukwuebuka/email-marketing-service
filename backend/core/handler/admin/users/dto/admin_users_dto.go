package dto

import "time"

type AdminFetchUserDTO struct {
	Search    string `json:"search"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	CompanyID string `json:"company_id"`
	UserID    string `json:"user_id"`
}

type AdminUserResponseDTO struct {
	ID                   string     `json:"id"`
	Fullname             string     `json:"fullname"`
	Email                string     `json:"email"`
	Phonenumber          *string    `json:"phonenumber"`
	Picture              *string    `json:"picture"`
	Verified             bool       `json:"verified"`
	Blocked              bool       `json:"blocked"`
	VerifiedAt           *time.Time `json:"verified_at"`
	Status               string     `json:"status"`
	ScheduledForDeletion bool       `json:"scheduled_for_deletion"`
	ScheduledDeletionAt  *time.Time `json:"scheduled_deletion_at"`
	LastLoginAt          *time.Time `json:"last_login_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	CompanyID            *string    `json:"company_id"`
	Companyname          *string    `json:"companyname"`
}


type UserResponse struct {
	ID                   string          `json:"id"`
	Fullname             string          `json:"fullname"`
	CompanyID            string          `json:"company_id"`
	Email                string          `json:"email"`
	Phonenumber          *string         `json:"phonenumber"`
	Password             *string         `json:"password"`
	GoogleID             *string         `json:"google_id"`
	Picture              *string         `json:"picture"`
	Verified             bool            `json:"verified"`
	Blocked              bool            `json:"blocked"`
	VerifiedAt           *time.Time      `json:"verified_at"`
	Status               string          `json:"status"`
	ScheduledForDeletion bool            `json:"scheduled_for_deletion"`
	ScheduledDeletionAt  *time.Time      `json:"scheduled_deletion_at"`
	LastLoginAt          *time.Time      `json:"last_login_at"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	DeletedAt            *time.Time      `json:"deleted_at"`
	Company              CompanyResponse `json:"company"`
}

type CompanyResponse struct {
	CompanyID        string     `json:"company_id"`
	Companyname      *string    `json:"companyname"`
	CompanyCreatedAt time.Time  `json:"company_created_at"`
	CompanyUpdatedAt time.Time  `json:"company_updated_at"`
	CompanyDeletedAt *time.Time `json:"company_deleted_at"`
}

 type AdminEmailLogDTO struct {
	Subject string `json:"subject" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type AdminUserStats struct {
	TotalContacts      int64 `json:"total_contacts"`
	TotalCampaigns     int64 `json:"total_campaigns"`
	TotalTemplates     int64 `json:"total_templates"`
	TotalCampaignsSent int64 `json:"total_campaigns_sent"`
	TotalGroups        int64 `json:"total_groups"`
}
