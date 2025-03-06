package requests

import "time"

// UpdateImage struct to update the image.
type UpdateImage struct {
	Name           *string   `json:"name"`
	Data           *string   `json:"data"`
	Description    *string   `json:"description"`
	UpdatedAt      time.Time `json:"updatedAt" validate:"required"`
	Quality        *int      `json:"quality"`
	IsNotResizable *bool     `json:"isNotResizable"`
}
