package responses

import "api-file/main/src/models"

// AppStoragePathPaginate struct for the AppStoragePath response.
type AppStoragePathPaginate struct {
	ID      uint   `json:"id"`
	AppName string `json:"app_name"`
	Path    string `json:"path"`
	Limit   *int64 `json:"limit"`
}

// SetAppStoragePathPaginate sets the AppStoragePath response.
func (response *AppStoragePathPaginate) SetAppStoragePathPaginate(appStoragePath *models.AppStoragePath) {
	response.ID = appStoragePath.ID
	response.AppName = appStoragePath.AppName
	response.Path = appStoragePath.Path

	if appStoragePath.Limit.Valid {
		response.Limit = &appStoragePath.Limit.Int64
	}
}
