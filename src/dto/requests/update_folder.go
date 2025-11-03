package requests

import "time"

// UpdateFolder struct for the PUT request.
type UpdateFolder struct {
	AppStoragePathID uint   `json:"appStoragePathId" validate:"required"`
	ParentFolderID   *uint  `json:"parentFolderId"`
	Name             string `json:"name" validate:"required"`
	// TODO: Add color validation.
	Color     string    `json:"color"`
	Immutable bool      `json:"immutable"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}
