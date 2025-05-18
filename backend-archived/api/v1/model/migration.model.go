package model 

import "time"

// MigrationHistory keeps track of migrations
type MigrationHistory struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"uniqueIndex"`
	AppliedAt time.Time
}