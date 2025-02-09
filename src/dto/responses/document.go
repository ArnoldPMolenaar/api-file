package responses

import (
	"api-file/main/src/models"
	"time"
)

type Document struct {
	ID        uint      `json:"id"`
	FolderID  uint      `json:"folderId"`
	Name      string    `json:"name"`
	Extension string    `json:"extension"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (d *Document) SetDocument(document *models.Document) {
	d.ID = document.ID
	d.FolderID = document.FolderID
	d.Name = document.Name
	d.Extension = document.Extension
	d.Size = document.Size
	d.CreatedAt = document.CreatedAt
	d.UpdatedAt = document.UpdatedAt
}
