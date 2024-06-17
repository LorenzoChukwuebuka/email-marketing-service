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

type SupportTicket struct {
	ID          uint          `gorm:"primaryKey" json:"-"`
	UUID        string        `gorm:"type:uuid;default:uuid_generate_v4();index"`
	Subject     string        `json:"subject"`
	Description string        `json:"description"`
	Status      Status        `json:"status"`
	Priority    Priority      `json:"priority"`
	SenderID    uint          `json:"sender_id"`
	AssignedTo  uint          `json:"assigned_to"`
	CreatedAt   time.Time     `json:"created_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	Files       []TicketFiles `gorm:"foreignKey:TicketID" json:"ticket_files"`
}

type TicketFiles struct {
	ID         uint      `gorm:"primary_key" json:"-"`
	UUID       string    `gorm:"type:uuid;default:uuid_generate_v4();index"`
	TicketID   uint      `json:"ticket_id"`
	TicketFile string    `json:"ticket_file"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// KnowledgeBaseArticle represents an article in the knowledge base
// type KnowledgeBaseArticle struct {
// 	ID         uint      `gorm:"primary_key" json:"id"`
// 	Title      string    `json:"title"`
// 	Content    string    `gorm:"type:text" json:"content"`
// 	CategoryID uint      `json:"category_id"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// }

// // KnowledgeBaseCategory represents a category for knowledge base articles
// type KnowledgeBaseCategory struct {
// 	ID          uint                   `gorm:"primary_key" json:"id"`
// 	Name        string                 `json:"name"`
// 	Description string                 `json:"description"`
// 	ParentID    *uint                  `json:"parent_id"`
// 	CreatedAt   time.Time              `json:"created_at"`
// 	UpdatedAt   time.Time              `json:"updated_at"`
// 	Articles    []KnowledgeBaseArticle `gorm:"foreignkey:CategoryID" json:"articles"`
// }
