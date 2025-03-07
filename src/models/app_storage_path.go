package models

import "database/sql"

type AppStoragePath struct {
	ID      uint   `gorm:"primaryKey"`
	AppName string `gorm:"not null;index:idx_app_storage_path,unique,priority:1"`
	Path    string `gorm:"not null;index:idx_app_storage_path,unique,priority:2"`
	Limit   sql.NullInt64

	// Relationships.
	App     App      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
	Folders []Folder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppStoragePathID;references:ID"`
}
