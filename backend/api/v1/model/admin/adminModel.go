package adminmodel

import (
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID         int       `json:"-" gorm:"primaryKey"`
	UUID       string    `json:"uuid"`
	FirstName  *string   `json:"firstname"`
	MiddleName *string   `json:"middlename"`
	LastName   *string   `json:"lastname"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt  time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type AdminResponse struct {
	ID         int     `json:"-"`
	UUID       string  `json:"uuid"`
	FirstName  *string `json:"firstname"`
	MiddleName *string `json:"middlename"`
	LastName   *string `json:"lastname"`
	Email      string  `json:"email"`
	Password   string  `json:"-"`
	Type       string  `jsosn:"type"`
}

type AdminChangePassword struct {
	AdminId     int    `json:"admin_id" `
	OldPassword []byte `json:"old_password"`
	NewPassword []byte `json:"new_password"`
}

type AdminMailLog struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	Mail string `json:"mail"`
}

type SystemsSMTPSetting struct {
	gorm.Model
	UUID           string `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	TXTRecord      string `json:"txt_record"`
	DMARCRecord    string `json:"dmarc_record"`
	DKIMSelector   string `json:"dkim_selector"`
	DKIMPublicKey  string `json:"dkim_public_key"`
	DKIMPrivateKey string `json:"dkim_private_key"`
	SPFRecord      string `json:"spf_record"`
	Verified       bool   `json:"verified"`
	MXRecord       string `gorm:"type:text"`
	Domain         string `json:"domain"`
}
