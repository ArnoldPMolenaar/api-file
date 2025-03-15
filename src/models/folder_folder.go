package models

type FolderFolder struct {
	AppStoragePathID uint `gorm:"primaryKey;autoIncrement:false"`
	FolderID         uint `gorm:"primaryKey;autoIncrement:false"`
	ParentFolderID   uint `gorm:"primaryKey;autoIncrement:false"`

	// Relationships.
	AppStoragePath AppStoragePath `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppStoragePathID;references:ID"`
	Folder         Folder         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FolderID;references:ID"`
	ParentFolder   Folder         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ParentFolderID;references:ID"`
}
