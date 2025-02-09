package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	FolderID  uint   `gorm:"not null"`
	Name      string `gorm:"not null"`
	Extension string `gorm:"not null"`
	Size      int    `gorm:"not null"`

	// Relationships.
	Folder Folder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID;references:ID"`
}
