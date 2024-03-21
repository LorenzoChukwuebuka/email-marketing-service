package model

import (
    "time"
)

// Ticket represents a support ticket
type Ticket struct {
    ID          uint `gorm:"primary_key"`
    Subject     string
    Description string
    Status      string `gorm:"type:enum('open', 'pending', 'resolved', 'closed')"`
    Priority    string `gorm:"type:enum('low', 'medium', 'high')"`
    ContactID   uint
    AssignedTo  uint
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Messages    []TicketMessage
}

// TicketMessage represents a message in a support ticket
type TicketMessage struct {
    ID        uint `gorm:"primary_key"`
    TicketID  uint
    Message   string
    Sender    string `gorm:"type:enum('user', 'agent')"`
    SenderID  uint
    CreatedAt time.Time
}

// KnowledgeBaseArticle represents an article in the knowledge base
type KnowledgeBaseArticle struct {
    ID         uint `gorm:"primary_key"`
    Title      string
    Content    string `gorm:"type:longtext"`
    CategoryID uint
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// KnowledgeBaseCategory represents a category for knowledge base articles
type KnowledgeBaseCategory struct {
    ID          uint `gorm:"primary_key"`
    Name        string
    Description string
    ParentID    *uint
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Articles    []KnowledgeBaseArticle `gorm:"foreignkey:CategoryID"`
}

