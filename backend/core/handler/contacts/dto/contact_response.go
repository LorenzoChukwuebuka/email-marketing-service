package dto

import (
	"github.com/google/uuid"
	"time"
)

type ContactResponse struct {
	ContactID        uuid.UUID         `json:"contact_id"`
	CompanyID        uuid.UUID         `json:"company_id"`
	FirstName        string            `json:"first_name"`
	LastName         string            `json:"last_name"`
	Email            string            `json:"email"`
	FromOrigin       string            `json:"from_origin"`
	IsSubscribed     bool              `json:"is_subscribed"`
	UserID           uuid.UUID         `json:"user_id"`
	ContactCreatedAt *time.Time        `json:"contact_created_at"`
	ContactUpdatedAt *time.Time        `json:"contact_updated_at"`
	UserContactGroup *UserContactGroup `json:"user_contact_group"`
	Group            *GroupInfo        `json:"group"`
}

type UserContactGroup struct {
	ID             *uuid.UUID `json:"id"`
	ContactID      *uuid.UUID `json:"contact_id"`
	ContactGroupID *uuid.UUID `json:"contact_group_id"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type GroupInfo struct {
	GroupID        *uuid.UUID `json:"group_id"`
	GroupName      string     `json:"group_name"`
	Description    *string    `json:"description"`
	GroupCreatorID *uuid.UUID `json:"group_creator_id"`
	GroupCreatedAt *time.Time `json:"group_created_at"`
	GroupUpdatedAt *time.Time `json:"group_updated_at"`
}

type GroupwithContactResponse struct {
	GroupID        uuid.UUID         `json:"group_id"`
	GroupName      string            `json:"group_name"`
	Description    *string           `json:"description"`
	GroupCreatedAt *time.Time        `json:"group_created_at"`
	Contacts       []GroupContactResponse `json:"contacts"`
}

type GroupContactResponse struct {
	ContactID           uuid.NullUUID `json:"contact_id"`
	ContactFirstName    *string       `json:"contact_first_name"`
	ContactLastName     *string       `json:"contact_last_name"`
	ContactEmail        *string       `json:"contact_email"`
	ContactFromOrigin   *string       `json:"contact_from_origin"`
	ContactIsSubscribed *bool         `json:"contact_is_subscribed"`
	ContactCreatedAt    *time.Time    `json:"contact_created_at"`
}
