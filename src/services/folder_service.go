package services

import (
	"api-file/main/src/database"
	"api-file/main/src/models"
)

// IsFolderAvailable method to check if a folder already exists inside the same path.
func IsFolderAvailable(appStoragePathID uint, folder string, parentFolderId ...uint) (bool, error) {
	if len(parentFolderId) != 0 {
		var folderIDs []int
		if result := database.Pg.Model(&models.FolderFolder{}).
			Where("parent_folder_id = ?", parentFolderId[0]).
			Pluck("folder_id", &folderIDs); result.Error != nil {
			return false, result.Error
		} else if len(folderIDs) == 0 {
			return false, nil
		}

		result := database.Pg.Limit(1).Find(&models.Folder{}, "name = ? AND id IN (?)", folder, folderIDs)
		if result.Error != nil {
			return false, result.Error
		} else {
			return result.RowsAffected == 1, nil
		}
	}

	if result := database.Pg.Table("folders").
		Select("folders.id").
		Joins("LEFT JOIN folder_folders ON folders.id = folder_folders.folder_id").
		Where("folders.app_storage_path_id = ? AND folders.name = ? AND folder_folders.folder_id IS NULL", appStoragePathID, folder).
		Limit(1).
		Find(&models.Folder{}); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsFolderDeleted method to check if a folder is deleted.
func IsFolderDeleted(id uint) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Folder{}).
		Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// GetFolderPath method to get the path of a folder.
// It returns the path of the folder like:
//
//	folder1/folder2/folder3
func GetFolderPath(appStoragePathID, folderID uint) (string, error) {
	var folders []*models.FolderFolder
	parentFolderID := folderID
	var path string

	if result := database.Pg.Preload("Folder").
		Preload("ParentFolder").
		Find(&folders, "app_storage_path_id = ?", appStoragePathID); result.Error != nil {
		return "", result.Error
	}

	folder := searchFolderByID(folders, folderID)
	for folder != nil {
		path = folder.Folder.Name + "/" + path
		parentFolderID = folder.ParentFolderID
		folder = searchFolderByID(folders, folder.ParentFolderID)
	}

	if mainFolder, _, err := GetFolder(parentFolderID); err != nil {
		return "", err
	} else {
		path = mainFolder.Name + "/" + path
	}

	return path, nil
}

// GetFolder method to get a folder.
func GetFolder(id uint, preload ...bool) (*models.Folder, []*models.Folder, error) {
	var folders []*models.Folder
	folder := &models.Folder{}
	query := database.Pg

	if len(preload) > 0 && preload[0] {
		query = query.Preload("Folders").Preload("Images").Preload("Documents")
	}

	if result := query.Find(folder, "id = ?", id); result.Error != nil {
		return nil, folders, result.Error
	}

	if len(folder.Folders) > 0 {
		folderIDs := make([]uint, len(folder.Folders))
		for i, f := range folder.Folders {
			folderIDs[i] = f.FolderID
		}
		if result := database.Pg.Find(&folders, "id IN (?)", folderIDs); result.Error != nil {
			return nil, folders, result.Error
		}
	}

	return folder, folders, nil
}

// CreateFolder method to create a folder.
func CreateFolder(appStoragePathID uint, name, color string, parentFolderID ...uint) (*models.Folder, error) {
	folder := &models.Folder{AppStoragePathID: appStoragePathID, Name: name, Color: color}

	if result := database.Pg.Create(folder); result.Error != nil {
		return nil, result.Error
	}

	if len(parentFolderID) > 0 {
		folderFolder := &models.FolderFolder{AppStoragePathID: appStoragePathID, FolderID: folder.ID, ParentFolderID: parentFolderID[0]}
		if result := database.Pg.Create(folderFolder); result.Error != nil {
			return nil, result.Error
		}
	}

	return folder, nil
}

// UpdateFolder method to update a folder.
func UpdateFolder(oldFolder *models.Folder, name, color string) (*models.Folder, error) {
	oldFolder.Name = name
	oldFolder.Color = color

	if result := database.Pg.Save(oldFolder); result.Error != nil {
		return nil, result.Error
	}

	return oldFolder, nil
}

// DeleteFolder method to delete a folder.
func DeleteFolder(folder *models.Folder) error {
	if result := database.Pg.Delete(folder); result.Error != nil {
		return result.Error
	}

	return nil
}

// RestoreFolder method to restore a folder.
func RestoreFolder(id uint) error {
	if result := database.Pg.Model(&models.Folder{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil); result.Error != nil {
		return result.Error
	}

	return nil
}

// searchInFoldersByID searches for a folder in the array by FolderID.
func searchFolderByID(folders []*models.FolderFolder, id uint) *models.FolderFolder {
	for _, folder := range folders {
		if folder.FolderID == id {
			return folder
		}
	}
	return nil
}
