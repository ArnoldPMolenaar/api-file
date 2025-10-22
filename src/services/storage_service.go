package services

import (
	"api-file/main/src/database"
	"api-file/main/src/models"
	"database/sql"
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

// IsStorageSpaceAvailable method to check if there is space available in the storage path.
func IsStorageSpaceAvailable(appStoragePathID uint) (bool, error) {
	var limit sql.NullInt64
	usedSpace, err := GetUsedSpace(appStoragePathID)
	if err != nil {
		return false, err
	}

	if result := database.Pg.Model(&models.AppStoragePath{}).
		Select("limit").
		Find(&limit, "id = ?", appStoragePathID).
		Scan(&limit); result.Error != nil {
		return false, result.Error
	}

	return !limit.Valid || usedSpace < limit.Int64, nil
}

// GetStoragePathIDByApp method to get the storage path ID by app name.
func GetStoragePathIDByApp(app string) (*uint, error) {
	var storagePathID *uint

	if result := database.Pg.Model(&models.AppStoragePath{}).
		Select("id").
		Where("app_name = ?", app).
		Limit(1).
		Scan(&storagePathID); result.Error != nil {
		return nil, result.Error
	}

	return storagePathID, nil
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

// GetUsedSpace method to get the used space for the app.
func GetUsedSpace(appStoragePathID uint) (int64, error) {
	var imagesSize, documentsSize int64

	// Sum images size
	if result := database.Pg.Model(&models.Image{}).
		Joins("JOIN folders ON images.folder_id = folders.id").
		Where("folders.app_storage_path_id = ?", appStoragePathID).
		Select("COALESCE(SUM(images.size), 0)").Scan(&imagesSize); result.Error != nil {
		return 0, result.Error
	}

	// Sum documents size
	if result := database.Pg.Model(&models.Document{}).
		Joins("JOIN folders ON documents.folder_id = folders.id").
		Where("folders.app_storage_path_id = ?", appStoragePathID).
		Select("COALESCE(SUM(documents.size), 0)").Scan(&documentsSize); result.Error != nil {
		return 0, result.Error
	}

	return imagesSize + documentsSize, nil
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
func CreateStoragePath(app, path string, limit *int64) (*models.AppStoragePath, error) {
	nullableLimit := sql.NullInt64{}
	if limit != nil {
		nullableLimit.Int64 = *limit
		nullableLimit.Valid = true
	} else {
		nullableLimit.Valid = false
	}

	storagePath := &models.AppStoragePath{AppName: app, Path: path, Limit: nullableLimit}

	if result := database.Pg.Create(storagePath); result.Error != nil {
		return nil, result.Error
	}

	return storagePath, nil
}

// UpdateStoragePath method to update a storage path for the app.
func UpdateStoragePath(oldStoragePath *models.AppStoragePath, app, path string, limit *int64) (*models.AppStoragePath, error) {
	oldStoragePath.AppName = app
	oldStoragePath.Path = path

	if limit != nil {
		oldStoragePath.Limit.Int64 = *limit
		oldStoragePath.Limit.Valid = true
	} else {
		oldStoragePath.Limit.Valid = false
	}

	if result := database.Pg.Save(oldStoragePath); result.Error != nil {
		return nil, result.Error
	}

	return oldStoragePath, nil
}
