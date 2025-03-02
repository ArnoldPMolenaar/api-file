package requests

// CreateDocument struct for creating a new document.
type CreateDocument struct {
	AppStoragePathID uint   `json:"appStoragePathId" validate:"required"`
	FolderID         uint   `json:"folderId" validate:"required"`
	Name             string `json:"name" validate:"required"`
	Data             string `json:"data" validate:"required"`
}
