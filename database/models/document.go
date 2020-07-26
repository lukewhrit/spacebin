package models

import (
	"github.com/jinzhu/gorm"
)

// Document is the structure of a document in the database
type Document struct {
	gorm.Model
	Key       string `gorm:"not null"`
	Content   string `gorm:"not null"`
	Extension string `gorm:"default: text"`
}
