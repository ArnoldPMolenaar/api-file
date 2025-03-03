package services

import (
	"api-file/main/src/database"
	"api-file/main/src/models"
)

// IsDocumentAvailable method to check if a document is available within the app.
func IsDocumentAvailable(folderId uint, name, extension string) (bool, error) {
	if result := database.Pg.
		Unscoped().
		Limit(1).
		Find(&models.Document{}, "folder_id = ? AND name = ? AND extension = ?", folderId, name, extension); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsDocumentDeleted method to check if a document is deleted.
func IsDocumentDeleted(id uint) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Document{}).
		Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// GetDocumentById method to get the document by its ID.
func GetDocumentById(id uint) (models.Document, error) {
	document := models.Document{}

	if result := database.Pg.
		Preload("Folder").
		Preload("Folder.AppStoragePath").
		Find(&document, "id = ?", id); result.Error != nil {
		return document, result.Error
	}

	return document, nil
}

// CreateDocument method to create a new document.
func CreateDocument(folderID uint, name, extension, mimeType string, size int) (models.Document, error) {
	document := models.Document{
		FolderID:  folderID,
		Name:      name,
		Extension: extension,
		MimeType:  mimeType,
		Size:      size,
	}

	if result := database.Pg.Create(&document); result.Error != nil {
		return models.Document{}, result.Error
	}

	return document, nil
}

// DeleteDocument method to delete a document.
func DeleteDocument(document *models.Document) error {
	if result := database.Pg.Delete(document); result.Error != nil {
		return result.Error
	}

	return nil
}

// RestoreDocument method to restore a document.
func RestoreDocument(id uint) error {
	if result := database.Pg.Model(&models.Document{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil); result.Error != nil {
		return result.Error
	}

	return nil
}
