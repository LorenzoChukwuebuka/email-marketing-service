package model

// Campaign represents an email campaign
// type Campaign struct {
//     ID           uint `gorm:"primary_key"`
//     Name         string
//     Subject      string
//     FromName     string
//     FromEmail    string
//     ReplyTo      string
//     Status       string `gorm:"type:enum('draft', 'scheduled', 'sending', 'sent', 'paused', 'cancelled')"`
//     CreatedAt    time.Time
//     UpdatedAt    time.Time
//     ScheduledAt  *time.Time
//     SentAt       *time.Time
//     Tracking     CampaignTracking
//     ListSegments []CampaignListSegment `gorm:"foreignkey:CampaignID"`
// }

// // CampaignTracking stores tracking data for a campaign
// type CampaignTracking struct {
//     ID              uint `gorm:"primary_key"`
//     CampaignID      uint
//     TotalRecipients int
//     TotalOpens      int
//     TotalClicks     int
//     UniqueOpens     int
//     UniqueClicks    int
//     Opens           []EmailOpen `gorm:"foreignkey:CampaignID"`
//     Clicks          []EmailClick `gorm:"foreignkey:CampaignID"`
// }

// // EmailOpen represents an email open event
// type EmailOpen struct {
//     ID         uint `gorm:"primary_key"`
//     CampaignID uint
//     ContactID  uint
//     OpenTime   time.Time
//     IPAddress  string
//     UserAgent  string
// }

// // EmailClick represents an email click event
// type EmailClick struct {
//     ID         uint `gorm:"primary_key"`
//     CampaignID uint
//     ContactID  uint
//     ClickTime  time.Time
//     LinkURL    string
//     IPAddress  string
//     UserAgent  string
// }

// // CampaignListSegment maps a campaign to a list segment
// type CampaignListSegment struct {
//     ID         uint `gorm:"primary_key"`
//     CampaignID uint
//     ListID     uint
//     SegmentID  uint
// }

// // ContactList represents a mailing list
// type ContactList struct {
//     ID          uint `gorm:"primary_key"`
//     Name        string
//     Description string
//     CreatedAt   time.Time
//     UpdatedAt   time.Time
//     Segments    []ListSegment `gorm:"foreignkey:ListID"`
// }

// // ListSegment represents a segment of a mailing list
// type ListSegment struct {
//     ID         uint `gorm:"primary_key"`
//     ListID     uint
//     Name       string
//     Criteria   string
//     CreatedAt  time.Time
//     UpdatedAt  time.Time
// }
