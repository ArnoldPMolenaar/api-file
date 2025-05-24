package responses

import "api-file/main/src/models"

// AppStoragePath struct for the AppStoragePath response.
type AppStoragePath struct {
	ID      uint     `json:"id"`
	AppName string   `json:"appName"`
	Path    string   `json:"path"`
	Limit   *int64   `json:"limit"`
	Used    int64    `json:"used"`
	Folders []Folder `json:"folders"`
}

// SetAppStoragePath sets the AppStoragePath response.
func (response *AppStoragePath) SetAppStoragePath(appStoragePath *models.AppStoragePath, usedSpace int64) {
	response.ID = appStoragePath.ID
	response.AppName = appStoragePath.AppName
	response.Path = appStoragePath.Path

	if appStoragePath.Limit.Valid {
		response.Limit = &appStoragePath.Limit.Int64
	}

	response.Used = usedSpace
	response.Folders = make([]Folder, len(appStoragePath.Folders))

	for i := range appStoragePath.Folders {
		response.Folders[i] = Folder{}
		response.Folders[i].SetFolder(&appStoragePath.Folders[i])
	}
}
