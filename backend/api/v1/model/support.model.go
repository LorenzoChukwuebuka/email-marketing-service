package model

import (
	adminmodel "email-marketing-service/api/v1/model/admin"
	"time"

	"gorm.io/gorm"
)

type Priority string

type SupportStatus string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

const (
	OpenTicket     SupportStatus = "open"
	CloseTicket    SupportStatus = "closed"
	ResolvedTicket SupportStatus = "resolved"
	PendingTicket  SupportStatus = "pending"
)

type SupportTicket struct {
	gorm.Model
	UUID         string          `gorm:"type:uuid;default:uuid_generate_v4();index"`
	UserID       string          `json:"user_id"`
	Name         string          `json:"name" gorm:"size:40"`
	Email        string          `json:"email" gorm:"size:40"`
	Subject      string          `json:"subject" gorm:"size:255"`
	Description  string          `json:"description" gorm:"type:text;default:null"`
	TicketNumber string          `json:"ticket_number"`
	Status       SupportStatus   `json:"status" gorm:"default:pending"`
	Priority     Priority        `json:"priority" gorm:"default:low"`
	LastReply    *time.Time      `json:"last_reply"`
	Files        []TicketFile    `gorm:"foreignKey:TicketID" json:"ticket_files"`
	Messages     []TicketMessage `gorm:"foreignKey:TicketID" json:"messages"`
}

type TicketFile struct {
	gorm.Model
	UUID     string `gorm:"type:uuid;default:uuid_generate_v4();index"`
	TicketID uint   `json:"ticket_id"`
	FileName string `json:"file_name" gorm:"size:255"`
	FilePath string `json:"file_path" gorm:"size:255"`
}

type TicketMessage struct {
	gorm.Model
	UUID     string `gorm:"type:uuid;default:uuid_generate_v4();index"`
	TicketID uint   `json:"ticket_id"`
	UserID   string `json:"user_id"`
	Message  string `json:"message" gorm:"type:text"`
	IsAdmin  bool   `json:"is_admin" gorm:"default:false"`
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
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	//    []KnowledgeBaseArticle `gorm:"foreignkey:CategoryID" json:"articles"`
}

type SupportTicketResponse struct {
	ID           uint                    `json:"-"`
	UUID         string                  `json:"uuid"`
	UserID       string                  `json:"user_id"`
	Name         string                  `json:"name"`
	Email        string                  `json:"email"`
	Subject      string                  `json:"subject"`
	Description  string                  `json:"description"`
	TicketNumber string                  `json:"ticket_number"`
	Status       SupportStatus           `json:"status"`
	Priority     Priority                `json:"priority"`
	LastReply    *time.Time              `json:"last_reply"`
	Files        []TicketFileResponse    `json:"files"`
	Messages     []TicketMessageResponse `json:"messages"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type TicketFileResponse struct {
	ID        uint      `json:"-"`
	UUID      string    `json:"uuid"`
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
}

type TicketMessageResponse struct {
	ID        uint                      `json:"-"`
	UUID      string                    `json:"uuid"`
	UserID    string                    `json:"user_id"`
	Message   string                    `json:"message"`
	IsAdmin   bool                      `json:"is_admin"`
	CreatedAt time.Time                 `json:"created_at"`
	User      *UserResponse             `json:"user,omitempty"`
	Admin     *adminmodel.AdminResponse `json:"admin,omitempty"`
}
