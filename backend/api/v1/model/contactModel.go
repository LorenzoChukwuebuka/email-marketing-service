package model

import (
	"gorm.io/gorm"
)

// Contact represents a contact entity.
type Contact struct {
	gorm.Model
	UUID         string         `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Email        string         `json:"email"`
	From         string         `json:"from"`
	IsSubscribed bool           `json:"is_subscribed"`
	UserId       string         `json:"user_id" gorm:"type:uuid"`
	Groups       []ContactGroup `json:"groups" gorm:"many2many:user_contact_groups;"`
}

type UserContactGroup struct {
	gorm.Model
	UUID           string       `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserId         string       `json:"user_id" gorm:"type:uuid"`
	ContactGroupId uint         `json:"contact_group_id" gorm:"column:contact_group_id"`
	ContactId      uint         `json:"contact_id"`
	Group          ContactGroup `json:"-" gorm:"foreignKey:ContactGroupId"`
	Contact        Contact      `json:"-" gorm:"foreignKey:ContactId"`
}

type ContactGroup struct {
	gorm.Model
	UUID        string    `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
	GroupName   string    `json:"group_name"`
	UserId      string    `json:"user_id"`
	Description string    `json:"description"`
	Contacts    []Contact `json:"contacts" gorm:"many2many:user_contact_groups;"`
}

type ContactResponse struct {
	ID        uint                   `json:"-"`
	UUID      string                 `json:"uuid"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Email     string                 `json:"email"`
	From      string                 `json:"from"`
	UserId    string                 `json:"user_id"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	DeletedAt *string                `json:"deleted_at"`
	Groups    []ContactGroupResponse `json:"groups"`
}

type ContactGroupResponse struct {
	ID          uint              `json:"-"`
	UUID        string            `json:"uuid"`
	GroupName   string            `json:"group_name"`
	UserId      string            `json:"user_id"`
	Description string            `json:"description"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	DeletedAt   *string           `json:"deleted_at"`
	Contacts    []ContactResponse `json:"contacts"`
}
