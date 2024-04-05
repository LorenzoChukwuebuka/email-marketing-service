package model

import (
	"time"
)

type Status string

const (
	Open     Status = "open"
	Pending  Status = "pending"
	Resolved Status = "resolved"
	Closed   Status = "closed"
)

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

// Ticket represents a support ticket
type Ticket struct {
	ID          uint            `gorm:"primary_key" json:"id"`
	Subject     string          `json:"subject"`
	Description string          `json:"description"`
	Status      Status          ` json:"status"`
	Priority    Priority        ` json:"priority"`
	ContactID   uint            `json:"contact_id"`
	AssignedTo  uint            `json:"assigned_to"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Messages    []TicketMessage `json:"messages"`
	Files       []TicketFiles   `json:"ticket_files"`
}

// TicketMessage represents a message in a support ticket
type TicketMessage struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	TicketID  uint      `json:"ticket_id"`
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
	SenderID  uint      `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
}

//ticket files

type TicketFiles struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	TicketID   uint      `json:"ticket_id"`
	TicketFile string    `json:"ticket_file"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// KnowledgeBaseArticle represents an article in the knowledge base
type KnowledgeBaseArticle struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Title      string    `json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// KnowledgeBaseCategory represents a category for knowledge base articles
type KnowledgeBaseCategory struct {
	ID          uint                   `gorm:"primary_key" json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ParentID    *uint                  `json:"parent_id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Articles    []KnowledgeBaseArticle `gorm:"foreignkey:CategoryID" json:"articles"`
}
