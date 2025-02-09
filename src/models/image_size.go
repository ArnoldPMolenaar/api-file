package models

import (
	"api-file/main/src/enums"
	"gorm.io/gorm"
)

type ImageSize struct {
	gorm.Model
	ImageID uint       `gorm:"not null"`
	Size    enums.Size `gorm:"not null;type:size"`
	Width   int        `gorm:"not null"`
	Height  int        `gorm:"not null"`

	// Relationships.
	Image Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ImageID;references:ID"`
}
