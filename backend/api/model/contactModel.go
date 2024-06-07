package model

import (
    "time"
)

// Contact represents a contact entity.
type Contact struct {
    ID            uint              `gorm:"primaryKey" json:"id"`
    FirstName     string            `json:"first_name"`
    LastName      string            `json:"last_name"`
    Company       string            `json:"company"`
    JobTitle      string            `json:"job_title"`
    DateOfBirth   time.Time         `json:"date_of_birth"`
    CreatedAt     time.Time         `json:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at"`
    Addresses     []Address         `json:"addresses"`
    Emails        []Email           `json:"emails"`
    Phones        []Phone           `json:"phones"`
    SocialMedia   []SocialMedia     `json:"social_media"`
    Notes         []Note            `json:"notes"`
    Tags          []ContactTag      `gorm:"many2many:contact_tag" json:"tags"`
    CustomFields  []ContactCustomField `json:"custom_fields"`
}

// Address represents the contact's address.
type Address struct {
    ID            uint       `gorm:"primaryKey" json:"id"`
    ContactID     uint       `json:"contact_id"`
    AddressType   string     `json:"address_type"`
    StreetAddress string     `json:"street_address"`
    City          string     `json:"city"`
    State         string     `json:"state"`
    PostalCode    string     `json:"postal_code"`
    Country       string     `json:"country"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
}

// Email represents the contact's email.
type Email struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    ContactID    uint       `json:"contact_id"`
    EmailAddress string     `json:"email_address"`
    EmailType    string     `json:"email_type"`
    IsPrimary    bool       `json:"is_primary"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// Phone represents the contact's phone number.
type Phone struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    ContactID    uint       `json:"contact_id"`
    PhoneNumber  string     `json:"phone_number"`
    PhoneType    string     `json:"phone_type"`
    IsPrimary    bool       `json:"is_primary"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// SocialMedia represents the contact's social media profile.
type SocialMedia struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    ContactID    uint       `json:"contact_id"`
    Platform     string     `json:"platform"`
    ProfileURL   string     `json:"profile_url"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// Note represents the note associated with a contact.
type Note struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    ContactID    uint       `json:"contact_id"`
    NoteContent  string     `json:"note_content"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// ContactTag represents the tag associated with a contact.
type ContactTag struct {
    ID        uint         `gorm:"primaryKey" json:"id"`
    TagName   string       `json:"tag_name"`
    Contacts  []Contact    `gorm:"many2many:contact_tag" json:"contacts"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
}

// ContactCustomField represents custom fields associated with a contact.
type ContactCustomField struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    ContactID   uint       `json:"contact_id"`
    FieldName   string     `json:"field_name"`
    FieldValue  string     `json:"field_value"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
