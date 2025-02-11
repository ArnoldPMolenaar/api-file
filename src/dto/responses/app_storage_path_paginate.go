package responses

import "api-file/main/src/models"

// AppStoragePathPaginate struct for the AppStoragePath response.
type AppStoragePathPaginate struct {
	ID      uint   `json:"id"`
	AppName string `json:"app_name"`
	Path    string `json:"path"`
}

// SetAppStoragePathPaginate sets the AppStoragePath response.
func (response *AppStoragePathPaginate) SetAppStoragePathPaginate(appStoragePath *models.AppStoragePath) {
	response.ID = appStoragePath.ID
	response.AppName = appStoragePath.AppName
	response.Path = appStoragePath.Path
}
