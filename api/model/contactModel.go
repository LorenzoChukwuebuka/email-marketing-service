package model

import (
    "time"
)


type Contact struct {
    ID        uint `gorm:"primaryKey"`
    FirstName string
    LastName  string
    Company   string
    JobTitle  string
    DateOfBirth time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
    Addresses   []Address
    Emails      []Email
    Phones      []Phone
    SocialMedia []SocialMedia
    Notes       []Note
    Tags        []ContactTag `gorm:"many2many:contact_tag"`
    CustomFields []ContactCustomField
}

// Address represents the contact's address.
type Address struct {
    ID           uint `gorm:"primaryKey"`
    ContactID    uint
    AddressType  string
    StreetAddress string
    City         string
    State        string
    PostalCode   string
    Country      string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Email represents the contact's email.
type Email struct {
    ID         uint `gorm:"primaryKey"`
    ContactID  uint
    EmailAddress string
    EmailType    string
    IsPrimary   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Phone represents the contact's phone number.
type Phone struct {
    ID          uint `gorm:"primaryKey"`
    ContactID   uint
    PhoneNumber string
    PhoneType   string
    IsPrimary  bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// SocialMedia represents the contact's social media profile.
type SocialMedia struct {
    ID         uint `gorm:"primaryKey"`
    ContactID  uint
    Platform   string
    ProfileURL string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// Note represents the note associated with a contact.
type Note struct {
    ID         uint `gorm:"primaryKey"`
    ContactID  uint
    NoteContent string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// ContactTag represents the tag associated with a contact.
type ContactTag struct {
    ID        uint `gorm:"primaryKey"`
    TagName   string
    Contacts  []Contact `gorm:"many2many:contact_tag"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// ContactCustomField represents custom fields associated with a contact.
type ContactCustomField struct {
    ID          uint `gorm:"primaryKey"`
    ContactID   uint
    FieldName   string
    FieldValue  string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
