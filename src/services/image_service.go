package services

import (
	"api-file/main/src/database"
	"api-file/main/src/models"
	"database/sql"
)

// IsImageAvailable method to check if an image is available within the app.
func IsImageAvailable(folderId uint, name, extension string) (bool, error) {
	if result := database.Pg.
		Limit(1).
		Find(&models.Image{}, "folder_id = ? AND name = ? AND extension = ?", folderId, name, extension); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// CreateImage method to create the image that is uploaded.
func CreateImage(folderID uint, name, extension, mimeType string, size, with, height int, description *string, sizes []models.ImageSize) (models.Image, error) {
	image := models.Image{
		FolderID:    folderID,
		Name:        name,
		Extension:   extension,
		MimeType:    mimeType,
		Size:        size,
		Width:       with,
		Height:      height,
		Description: sql.NullString{Valid: false, String: ""},
		ImageSizes:  sizes,
	}

	if description != nil {
		image.Description.Valid = true
		image.Description.String = *description
	}

	if result := database.Pg.Create(&image); result.Error != nil {
		return models.Image{}, result.Error
	}

	return image, nil
}
