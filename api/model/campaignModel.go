package model

import (
    "time"
)

// Campaign represents an email campaign
type Campaign struct {
    ID           uint `gorm:"primary_key"`
    Name         string
    Subject      string
    FromName     string
    FromEmail    string
    ReplyTo      string
    Status       string `gorm:"type:enum('draft', 'scheduled', 'sending', 'sent', 'paused', 'cancelled')"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    ScheduledAt  *time.Time
    SentAt       *time.Time
    Lists        []CampaignList
    Segments     []CampaignSegment
    Content      CampaignContent
    Tracking     CampaignTracking
    Clicks       []EmailClick
    Opens        []EmailOpen
    Links        []CampaignLink
}

// CampaignList maps a campaign to a contact list
type CampaignList struct {
    ID         uint `gorm:"primary_key"`
    CampaignID uint
    ListID     uint
}

// CampaignSegment represents a segment for a campaign
type CampaignSegment struct {
    ID             uint `gorm:"primary_key"`
    CampaignID     uint
    SegmentName    string
    SegmentCriteria string
}

// CampaignContent stores the email content for a campaign
type CampaignContent struct {
    ID           uint `gorm:"primary_key"`
    CampaignID   uint
    HTMLContent  string `gorm:"type:longtext"`
    TextContent  string `gorm:"type:longtext"`
}

// CampaignTracking stores tracking data for a campaign
type CampaignTracking struct {
    ID              uint `gorm:"primary_key"`
    CampaignID      uint
    TotalRecipients int
    TotalOpens      int
    TotalClicks     int
}

// EmailOpen represents an email open event
type EmailOpen struct {
    ID         uint `gorm:"primary_key"`
    CampaignID uint
    ContactID  uint
    OpenTime   time.Time
    IPAddress  string
    UserAgent  string
}

// EmailClick represents an email click event
type EmailClick struct {
    ID         uint `gorm:"primary_key"`
    CampaignID uint
    ContactID  uint
    ClickTime  time.Time
    LinkURL    string
    IPAddress  string
    UserAgent  string
}

// CampaignLink represents a link used in a campaign
type CampaignLink struct {
    ID         uint `gorm:"primary_key"`
    CampaignID uint
    LinkURL    string
}