package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	FolderID    uint   `gorm:"not null;index:idx_image,unique,priority:1"`
	Name        string `gorm:"not null;index:idx_image,unique,priority:2"`
	Extension   string `gorm:"not null;index:idx_image,unique,priority:3"`
	MimeType    string `gorm:"not null"`
	Size        int    `gorm:"not null"`
	Width       int    `gorm:"not null"`
	Height      int    `gorm:"not null"`
	Description sql.NullString

	// Relationships.
	Folder     Folder      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID;references:ID"`
	ImageSizes []ImageSize `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ImageID"`
}
