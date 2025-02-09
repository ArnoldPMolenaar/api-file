package requests

// UpdateFolder struct for the PUT request.
type UpdateFolder struct {
	AppStoragePathID uint   `json:"appStoragePathId" validate:"required"`
	ParentFolderID   *uint  `json:"parentFolderId"`
	Name             string `json:"name" validate:"required"`
	// TODO: Add color validation.
	Color string `json:"color"`
}
