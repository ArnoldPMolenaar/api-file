package requests

// UpdateAppStoragePath struct for updating an AppStoragePath record.
type UpdateAppStoragePath struct {
	App   string `json:"app" validate:"required"`
	Path  string `json:"path" validate:"required"`
	Limit *int64 `json:"limit"`
}
