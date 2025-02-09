package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	AppStoragePathID uint   `gorm:"not null"`
	Name             string `gorm:"not null"`
	Extension        string `gorm:"not null"`
	Size             int    `gorm:"not null"`
	Path             string `gorm:"not null"`

	// Relationships.
	AppStoragePath AppStoragePath `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppStoragePathID;references:ID"`
}
