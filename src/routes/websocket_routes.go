package routes

import (
	"api-file/main/src/controllers"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func WebSocketRoutes(a *fiber.App) {
	route := a.Group("/v1")

	// Create websocket routes group.
	ws := route.Group("/ws")
	ws.Get("/progress", websocket.New(controllers.WebSocketProgress))
}
