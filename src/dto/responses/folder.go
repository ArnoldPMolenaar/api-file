package responses

import (
	"api-file/main/src/models"
	"time"
)

type Folder struct {
	ID               uint      `json:"id"`
	AppStoragePathID uint      `json:"appStoragePathId"`
	Name             string    `json:"name"`
	Color            string    `json:"color"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// SetFolder method to set a folder.
func (f *Folder) SetFolder(folder *models.Folder) {
	f.ID = folder.ID
	f.AppStoragePathID = folder.AppStoragePathID
	f.Name = folder.Name
	f.Color = folder.Color
	f.CreatedAt = folder.CreatedAt
	f.UpdatedAt = folder.UpdatedAt
}
