package routes

import (
	"api-file/main/src/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register CRUD routes for /v1/images.
	images := route.Group("/image")
	images.Get("/:id", controllers.GetImageFile)
	images.Get("/:id/:size", controllers.GetImageFileSize)
}
