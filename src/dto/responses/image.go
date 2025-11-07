package responses

import (
	"api-file/main/src/models"
	"time"
)

type Image struct {
	ID               uint        `json:"id"`
	FolderID         uint        `json:"folderId"`
	AppStoragePathID uint        `json:"appStoragePathId"`
	Name             string      `json:"name"`
	Extension        string      `json:"extension"`
	Size             int         `json:"size"`
	Width            int         `json:"width"`
	Height           int         `json:"height"`
	Description      *string     `json:"description"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
	ImageSizes       []ImageSize `json:"sizes"`
}

// SetImage method to set an image.
func (i *Image) SetImage(image *models.Image, appStoragePathID *uint) {
	i.ID = image.ID
	i.FolderID = image.FolderID
	i.Name = image.Name
	i.Extension = image.Extension
	i.Size = image.Size
	i.Width = image.Width
	i.Height = image.Height
	i.CreatedAt = image.CreatedAt
	i.UpdatedAt = image.UpdatedAt
	i.ImageSizes = []ImageSize{}

	if appStoragePathID != nil {
		i.AppStoragePathID = *appStoragePathID
	} else {
		i.AppStoragePathID = image.Folder.AppStoragePathID
	}

	if image.Description.Valid {
		i.Description = &image.Description.String
	}

	for index := range image.ImageSizes {
		imageSize := ImageSize{}
		imageSize.SetImageSize(&image.ImageSizes[index])
		i.ImageSizes = append(i.ImageSizes, imageSize)
	}
}
