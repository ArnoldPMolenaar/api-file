package routes

import (
	"api-file/main/src/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register CRUD routes for /v1/image.
	image := route.Group("/image")
	image.Get("/:id", controllers.GetImageFile)
	image.Get("/:id/:size", controllers.GetImageFileSize)

	// Register CRUD routes for /v1/document.
	document := route.Group("/document")
	document.Get("/:id", controllers.GetDocumentFile)
}
