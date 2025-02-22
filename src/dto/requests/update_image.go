package requests

import "time"

// UpdateImage struct to update the image.
type UpdateImage struct {
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt" validate:"required"`
}
