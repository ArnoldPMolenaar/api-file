package services

import (
	"api-file/main/src/database"
	"api-file/main/src/models"
	"database/sql"
)

// IsImageAvailable method to check if an image is available within the app.
func IsImageAvailable(folderId uint, name, extension string) (bool, error) {
	if result := database.Pg.
		Unscoped().
		Limit(1).
		Find(&models.Image{}, "folder_id = ? AND name = ? AND extension = ?", folderId, name, extension); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsImageDeleted method to check if a image is deleted.
func IsImageDeleted(id uint) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Image{}).
		Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// GetImageById method to get the image by its ID.
func GetImageById(id uint, withSizes bool) (models.Image, error) {
	image := models.Image{}
	query := database.Pg

	if withSizes {
		query = query.Preload("ImageSizes")
	}

	if result := query.Find(&image, "id = ?", id); result.Error != nil {
		return models.Image{}, result.Error
	}

	return image, nil
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

// UpdateImage method to update the image description.
func UpdateImage(image *models.Image, description string) (models.Image, error) {
	image.Description.Valid = description != ""
	image.Description.String = description

	if result := database.Pg.Save(&image); result.Error != nil {
		return models.Image{}, result.Error
	}

	return *image, nil
}

// DeleteImage method to delete a image.
func DeleteImage(image *models.Image) error {
	if result := database.Pg.Delete(image); result.Error != nil {
		return result.Error
	}

	return nil
}

// RestoreImage method to restore an image.
func RestoreImage(id uint) error {
	if result := database.Pg.Model(&models.Image{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil); result.Error != nil {
		return result.Error
	}

	return nil
}
