package requests

// CreateImage struct for creating a new image.
type CreateImage struct {
	AppStoragePathID uint    `json:"appStoragePathId" validate:"required"`
	FolderID         uint    `json:"folderId" validate:"required"`
	Name             string  `json:"name" validate:"required"`
	Data             string  `json:"data" validate:"required"`
	Description      *string `json:"description"`
	Quality          int     `json:"quality"`
	IsNotResizable   bool    `json:"isNotResizable"`
}
