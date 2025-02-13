package requests

// CreateImage struct for creating a new image.
type CreateImage struct {
	FolderID       uint    `json:"folderId" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	Data           string  `json:"data" validate:"required"`
	Description    *string `json:"description"`
	IsNotResizable bool    `json:"isNotResizable"`
}
