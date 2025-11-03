package models

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	AppStoragePathID uint   `gorm:"not null"`
	Name             string `gorm:"not null"`
	Color            string `gorm:"not null"`
	Immutable        bool   `gorm:"default:false;not null"`

	// Relationships.
	AppStoragePath AppStoragePath `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppStoragePathID;references:ID"`
	Folders        []FolderFolder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ParentFolderID"`
	Images         []Image        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID"`
	Documents      []Document     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID"`
}
