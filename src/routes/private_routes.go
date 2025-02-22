package routes

import (
	"api-file/main/src/controllers"
	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register CRU routes for /v1/storage-paths.
	storagePaths := route.Group("/storage-paths", middleware.MachineProtected())
	storagePaths.Get("/", controllers.GetStoragePaths)
	storagePaths.Post("/", controllers.CreateStoragePath)
	storagePaths.Get("/:id", controllers.GetStoragePath)
	storagePaths.Put("/:id", controllers.UpdateStoragePath)

	// Register CRUD routes for /v1/folders.
	folders := route.Group("/folders", middleware.MachineProtected())
	folders.Post("/", controllers.CreateFolder)
	folders.Get("/:id", controllers.GetFolder)
	folders.Put("/:id", controllers.UpdateFolder)
	folders.Delete("/:id", controllers.DeleteFolder)
	folders.Put("/:id/restore", controllers.RestoreFolder)

	// Register CRUD routes for /v1/images.
	images := route.Group("/images", middleware.MachineProtected())
	images.Post("/", controllers.CreateImage)
	images.Get("/:id", controllers.GetImage)
	images.Put("/:id", controllers.UpdateImage)
	images.Delete("/:id", controllers.DeleteImage)
	images.Put("/:id/restore", controllers.RestoreImage)

	// Register handshake route for websocket.
	route.Get("/handshake", middleware.MachineProtected(), controllers.Handshake)
}
