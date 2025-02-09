package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	FolderID    uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Extension   string `gorm:"not null"`
	Size        int    `gorm:"not null"`
	Width       int    `gorm:"not null"`
	Height      int    `gorm:"not null"`
	Description sql.NullString

	// Relationships.
	Folder     Folder      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID;references:ID"`
	ImageSizes []ImageSize `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ImageID"`
}
