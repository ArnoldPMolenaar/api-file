package responses

import (
	"api-file/main/src/models"
	"time"
)

type Document struct {
	ID               uint      `json:"id"`
	FolderID         uint      `json:"folderId"`
	AppStoragePathID uint      `json:"appStoragePathId"`
	Name             string    `json:"name"`
	Extension        string    `json:"extension"`
	Size             int       `json:"size"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// SetDocument sets the document properties.
func (d *Document) SetDocument(document *models.Document, appStoragePathID *uint) {
	d.ID = document.ID
	d.FolderID = document.FolderID
	d.Name = document.Name
	d.Extension = document.Extension
	d.Size = document.Size
	d.CreatedAt = document.CreatedAt
	d.UpdatedAt = document.UpdatedAt

	if appStoragePathID != nil {
		d.AppStoragePathID = *appStoragePathID
	} else {
		d.AppStoragePathID = document.Folder.AppStoragePathID
	}
}
