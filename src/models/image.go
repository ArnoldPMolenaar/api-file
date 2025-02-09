package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	AppStoragePathID uint   `gorm:"not null"`
	Name             string `gorm:"not null"`
	Extension        string `gorm:"not null"`
	Size             int    `gorm:"not null"`
	Path             string `gorm:"not null"`
	Width            int    `gorm:"not null"`
	Height           int    `gorm:"not null"`

	// Relationships.
	AppStoragePath AppStoragePath `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppStoragePathID;references:ID"`
	ImageSizes     []ImageSize    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ImageID"`
}
