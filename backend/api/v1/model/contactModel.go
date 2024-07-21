package model

import (
	"time"
)

// Contact represents a contact entity.
type Contact struct {
	ID        uint       `gorm:"primaryKey" json:"-"`
	UUID      string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	GroupId   uint       `json:"group_id" gorm:"default:null"`
	From      string     `json:"from"`
	UserId    string     `json:"user_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type ContactGroup struct {
	ID          uint       `gorm:"primaryKey" json:"-"`
	UUID        string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	GroupName   string     `json:"group_name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
	Contacts    []Contact  `json:"contacts" gorm:"foreignKey:GroupId"`
}

type ContactResponse struct {
	ID        uint   ` json:"-"`
	UUID      string `json:"uuid" `
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	GroupId   uint   `json:"group_id"`
	From      string `json:"from"`
	UserId    string `json:"user_id" `
	CreatedAt string `json:"created_at" `
	UpdatedAt string `json:"updated_at" `
	DeletedAt string `json:"deleted_at"`
}
