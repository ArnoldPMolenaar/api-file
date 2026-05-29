package requests

// CreateApp request DTO to create an App.
type CreateApp struct {
	Name string `json:"name" validate:"required"`
}
