package requests

// CreateAppStoragePath struct for creating a new AppStoragePath.
type CreateAppStoragePath struct {
	App  string `json:"app" validate:"required"`
	Path string `json:"path" validate:"required"`
}
