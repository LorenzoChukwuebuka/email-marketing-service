package model

import (
	"gorm.io/gorm"
	"time"
)

// Contact represents a contact entity.
type Contact struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	UUID      string         `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	From      string         `json:"from"`
	UserId    string         `json:"user_id" gorm:"type:uuid"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Groups    []ContactGroup `json:"groups" gorm:"many2many:user_contact_groups;"`
}

type UserContactGroup struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	UserId    string         `json:"user_id" gorm:"type:uuid"`
	GroupId   uint           `json:"group_id"`
	ContactId uint           `json:"contact_id"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Group     ContactGroup   `json:"-" gorm:"foreignKey:GroupId"`
	Contact   Contact        `json:"-" gorm:"foreignKey:ContactId"`
}

type ContactGroup struct {
	ID          uint           `gorm:"primaryKey" json:"-"`
	UUID        string         `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	GroupName   string         `json:"group_name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time     `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Contacts    []Contact      `json:"contacts" gorm:"many2many:user_contact_groups;"`
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
