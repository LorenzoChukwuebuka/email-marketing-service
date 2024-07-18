package model

import (
	"time"
)

// Contact represents a contact entity.
type Contact struct {
	ID        uint       `gorm:"primaryKey" json:"-"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	GroupId   uint       `json:"group_id" gorm:"default:null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
}

type ContactGroup struct {
	ID          uint       `gorm:"primaryKey" json:"-"`
	GroupName   string     `json:"group_name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
	Contacts    []Contact
}
