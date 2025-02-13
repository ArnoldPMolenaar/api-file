package routes

import (
	"api-file/main/src/controllers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebSocketRoutes(a *fiber.App) {
	// Create websocket routes group.
	route := a.Group("/ws")
	route.Get("/progress", websocket.New(controllers.WebSocketProgress))
}
