package requests

import "time"

// UpdateDocument struct for updating a document.
type UpdateDocument struct {
	Name      string    `json:"name" validate:"required"`
	Data      string    `json:"data" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}
