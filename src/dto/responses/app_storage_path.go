package responses

import "api-file/main/src/models"

// AppStoragePath struct for the AppStoragePath response.
type AppStoragePath struct {
	ID      uint   `json:"id"`
	AppName string `json:"app_name"`
	Path    string `json:"path"`
}

// SetAppStoragePath sets the AppStoragePath response.
func (response *AppStoragePath) SetAppStoragePath(appStoragePath *models.AppStoragePath) {
	response.ID = appStoragePath.ID
	response.AppName = appStoragePath.AppName
	response.Path = appStoragePath.Path
}
