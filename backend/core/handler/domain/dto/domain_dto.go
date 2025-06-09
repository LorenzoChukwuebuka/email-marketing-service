package dto

import "time"

type DomainDTO struct {
	Domain    string `json:"domain" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
	CompanyID string `json:"company_id" validate:"required"`
}

type FetchDomainDTO struct {
	CompanyID   string `json:"company_id"`
	UserID      string `json:"user_id"`
	DomainID    string `json:"domain_id"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	SearchQuery string `json:"search_query"`
}

type DomainResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	CompanyID      string     `json:"company_id"`
	Domain         string     `json:"domain"`
	TxtRecord      *string    `json:"txt_record"`
	DmarcRecord    *string    `json:"dmarc_record"`
	DkimSelector   *string    `json:"dkim_selector"`
	DkimPublicKey  *string    `json:"dkim_public_key"`
	DkimPrivateKey *string    `json:"dkim_private_key"`
	SpfRecord      *string    `json:"spf_record"`
	Verified       bool       `json:"verified"`
	MxRecord       *string    `json:"mx_record"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
