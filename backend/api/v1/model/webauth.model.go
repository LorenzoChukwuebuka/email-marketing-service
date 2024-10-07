package model

import (
	"gorm.io/gorm"
)

type WebAuthnCredential struct {
	gorm.Model
	UserID       string `gorm:"not null"`
	CredentialID []byte `gorm:"not null"`
	PublicKey    []byte `gorm:"not null"`
	SignCount    uint32 `gorm:"not null"`
}
