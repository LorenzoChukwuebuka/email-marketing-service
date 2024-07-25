package model

import (
	"time"
)

type EmailTemplate struct {
    
    Name         string
    Description  string
    HTMLContent  string `gorm:"type:longtext"`
    TextContent  string `gorm:"type:longtext"`
  
    Categories   []TemplateCategoryMapping `gorm:"foreignkey:TemplateID"`
}

// TemplateCategory represents a category for email templates
type TemplateCategory struct {
    ID          uint `gorm:"primary_key"`
    Name        string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Templates   []TemplateCategoryMapping `gorm:"foreignkey:CategoryID"`
}

// TemplateCategoryMapping maps email templates to categories
type TemplateCategoryMapping struct {
    ID         uint `gorm:"primary_key"`
    TemplateID uint
    CategoryID uint
}