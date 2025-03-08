package services

import (
	"api-file/main/src/cache"
	"api-file/main/src/database"
	"api-file/main/src/enums"
	"api-file/main/src/models"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
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
	query := database.Pg.Preload("Folder").Preload("Folder.AppStoragePath")

	if withSizes {
		query = query.Preload("ImageSizes")
	}

	if result := query.Find(&image, "id = ?", id); result.Error != nil {
		return models.Image{}, result.Error
	}

	return image, nil
}

// GetImageSizeById method to get the image size by its ImageID.
func GetImageSizeById(id uint, size enums.Size) (models.ImageSize, error) {
	imageSize := models.ImageSize{}

	if result := database.Pg.
		Preload("Image").
		Preload("Image.Folder").
		Preload("Image.Folder.AppStoragePath").
		Find(&imageSize, "image_id = ? AND size = ?", id, size); result.Error != nil {
		return imageSize, result.Error
	}

	return imageSize, nil
}

// GetImage method to get the image by its ID with preloading.
// Used by loading a file.
func GetImage(id uint) (models.Image, error) {
	image := models.Image{}

	if result := database.Pg.
		Preload("Folder").
		Preload("Folder.AppStoragePath").
		Find(&image, "id = ?", id); result.Error != nil {
		return models.Image{}, result.Error
	}

	return image, nil
}

// GetImageFromCache method to get the image from the cache.
func GetImageFromCache(id uint, size ...string) (string, error) {
	key := ImageCacheKey(id, size...)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return "", result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return "", err
	}

	return value, nil
}

// CreateImage method to create the image that is uploaded.
func CreateImage(folderID uint, name, extension, mimeType string, size, width, height int, description *string, sizes []models.ImageSize) (models.Image, error) {
	image := models.Image{
		FolderID:    folderID,
		Name:        name,
		Extension:   extension,
		MimeType:    mimeType,
		Size:        size,
		Width:       width,
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

// SaveImageToCache method to save the image to the cache.
func SaveImageToCache(imageId uint, path string, size ...string) error {
	key := ImageCacheKey(imageId, size...)

	expiration := os.Getenv("VALKEY_EXPIRATION_IMAGE")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(key).Value(path).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// UpdateImage method to update the image description.
func UpdateImage(image *models.Image, name, extension, mimeType *string, size, width, height *int, description *string, sizes *[]models.ImageSize) (models.Image, error) {
	if name != nil {
		image.Name = *name
	}
	if extension != nil {
		image.Extension = *extension
	}
	if mimeType != nil {
		image.MimeType = *mimeType
	}
	if size != nil {
		image.Size = *size
	}
	if width != nil {
		image.Width = *width
	}
	if height != nil {
		image.Height = *height
	}

	image.Description.Valid = description != nil && *description != ""
	if image.Description.Valid {
		image.Description.String = *description
	}

	if sizes != nil {
		if result := database.Pg.Model(&models.ImageSize{}).Unscoped().Delete(&models.ImageSize{}, "image_id = ?", image.ID); result.Error != nil {
			return *image, result.Error
		}
		image.ImageSizes = *sizes
	}

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

	if result := database.Pg.Model(&models.ImageSize{}).Delete(&models.ImageSize{}, "image_id = ?", image.ID); result.Error != nil {
		return result.Error
	}

	_ = DeleteImageFromCache(image.ID)
	for i := range image.ImageSizes {
		_ = DeleteImageFromCache(image.ID, image.ImageSizes[i].Size.String())
	}

	return nil
}

// DeleteImageFromCache method to delete the image from the cache.
func DeleteImageFromCache(id uint, size ...string) error {
	key := ImageCacheKey(id, size...)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
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

// ImageCacheKey method to create a cache key for the image.
func ImageCacheKey(id uint, size ...string) string {
	if len(size) > 0 {
		return fmt.Sprintf("image:%d:%s", id, size[0])
	}

	return fmt.Sprintf("image:%d", id)
}
