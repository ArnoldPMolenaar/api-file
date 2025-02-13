package middleware

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"strings"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App) {
	a.Use(
		// Add CORS to each route.
		cors.New(cors.Config{
			AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodPut,
				fiber.MethodHead,
				fiber.MethodOptions,
			}, ","),
			AllowHeaders: "Accept,Content-Type",
		}),

		// Add simple logger.
		logger.New(),

		// Catch a panic and return a 500 response.
		recover.New(),
	)
	a.Use("/ws", webSocketMiddleware)
}

// WebSocketMiddleware checks if the request is a WebSocket upgrade.
func webSocketMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
