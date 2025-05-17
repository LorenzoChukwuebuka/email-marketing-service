package model

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	UserID  string `gorm:"index"`
	Action  string `gorm:"type:varchar(255)"`
	Details string `gorm:"type:text"`
}
