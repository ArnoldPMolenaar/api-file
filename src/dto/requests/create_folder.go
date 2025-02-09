package requests

// CreateFolder struct for the POST request.
type CreateFolder struct {
	AppStoragePathID uint   `json:"appStoragePathId" validate:"required"`
	ParentFolderID   *uint  `json:"parentFolderId"`
	Name             string `json:"name" validate:"required"`
	// TODO: Add color validation.
	Color string `json:"color"`
}
