package model

import "gorm.io/gorm"

type Sender struct {
	gorm.Model
	UUID     string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	IsSigned bool   `json:"is_signed"`
	DomainID string `json:"domain_id"`
}

type Domains struct {
	gorm.Model
	UUID           string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID         string `json:"user_id"`
	Domain         string `json:"domain"`
	TXTRecord      string `json:"txt_record"`
	DMARCRecord    string `json:"dmarc_record"`
	DKIMSelector   string `json:"dkim_selector"`
	DKIMPublicKey  string `json:"dkim_public_key"`
	DKIMPrivateKey string `json:"dkim_private_key"`
	Verified       bool   `json:"verified"`
}

type SenderResponse struct {
	ID        uint    `json:"-"`
	UUID      string  `json:"uuid" `
	UserID    string  `json:"user_id" `
	Name      string  `json:"name" `
	Email     string  `json:"email" `
	Verified  bool    `json:"verified"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type DomainsResponse struct {
	ID            uint    `json:"-"`
	UUID          string  `json:"uuid" `
	UserID        string  `json:"user_id" `
	Domain        string  `json:"domain" `
	TXTRecord     string  `json:"txt_record"`
	DMARCRecord   string  `json:"dmarc_record"`
	DKIMSelector  string  `json:"dkim_selector"`
	DKIMPublicKey string  `json:"dkim_public_key"`
	Verified      bool    `json:"verified"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
}
