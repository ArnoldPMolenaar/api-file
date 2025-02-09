package models

type FolderFolder struct {
	AppStoragePathID uint `gorm:"not null"`
	FolderID         uint `gorm:"not null"`
	ParentFolderID   uint `gorm:"not null"`

	// Relationships.
	AppStoragePath AppStoragePath `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppStoragePathID;references:ID"`
	Folder         Folder         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID;references:ID"`
	ParentFolder   Folder         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ParentFolderID;references:ID"`
}
