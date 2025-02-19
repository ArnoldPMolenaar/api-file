package responses

import (
	"api-file/main/src/models"
	"time"
)

type ImageSize struct {
	ID        uint      `json:"id"`
	ImageID   uint      `json:"imageId"`
	Size      string    `json:"size"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetImageSize method to set an image size.
func (is *ImageSize) SetImageSize(imageSize *models.ImageSize) {
	is.ID = imageSize.ID
	is.ImageID = imageSize.ImageID
	is.Size = imageSize.Size.String()
	is.Width = imageSize.Width
	is.Height = imageSize.Height
	is.CreatedAt = imageSize.CreatedAt
	is.UpdatedAt = imageSize.UpdatedAt
}
