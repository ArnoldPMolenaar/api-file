package services

import (
	"api-file/main/src/database"
	"api-file/main/src/models"
	"os"
)

// IsStorageAvailable method to check if a storage path is available within the app.
func IsStorageAvailable(app, path string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.AppStoragePath{}, "app_name = ? AND path = ?", app, path); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetPath method to get the full path.
func GetPath(appStoragePath *models.AppStoragePath, folderID uint) (string, error) {
	path := os.Getenv("PATH_FILES") + appStoragePath.Path
	folderPath, err := GetFolderPath(appStoragePath.ID, folderID)

	if err != nil {
		return "", err
	}

	return path + folderPath, nil
}

// GetStoragePath method to get a storage path for the app.
func GetStoragePath(id uint) (*models.AppStoragePath, error) {
	storagePath := &models.AppStoragePath{}
	folders := make([]models.Folder, 0)

	if result := database.Pg.Joins("LEFT JOIN folder_folders ON folder_folders.folder_id = folders.id").
		Find(&folders, "folder_folders.folder_id IS NULL AND folders.app_storage_path_id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	if result := database.Pg.Find(storagePath, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	storagePath.Folders = folders

	return storagePath, nil
}

// CreateStoragePath method to create a storage path for the app.
func CreateStoragePath(app, path string) (*models.AppStoragePath, error) {
	storagePath := &models.AppStoragePath{AppName: app, Path: path}

	if result := database.Pg.Create(storagePath); result.Error != nil {
		return nil, result.Error
	}

	return storagePath, nil
}

// UpdateStoragePath method to update a storage path for the app.
func UpdateStoragePath(oldStoragePath *models.AppStoragePath, app, path string) (*models.AppStoragePath, error) {
	oldStoragePath.AppName = app
	oldStoragePath.Path = path

	if result := database.Pg.Save(oldStoragePath); result.Error != nil {
		return nil, result.Error
	}

	return oldStoragePath, nil
}
