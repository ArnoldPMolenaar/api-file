package controllers

import (
	"api-file/main/src/dto/responses"
	"api-file/main/src/errors"
	"api-file/main/src/services"
	"encoding/json"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

// ProgressConnections is a map of WebSocket connections from clients.
var ProgressConnections = make(map[*websocket.Conn]bool)

// WebSocketProgress is a WebSocket handler that sends progress updates to the client.
func WebSocketProgress(c *websocket.Conn) {
	// Access query parameters.
	app := c.Query("app")
	id := c.Query("id")
	code := c.Query("code")

	// Validate the handshake.
	if value, err := services.GetHandshake(app, code); err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Handshake failed: "+err.Error()))
		c.Close()
		return
	} else if value != id {
		c.WriteMessage(websocket.TextMessage, []byte("Handshake failed: invalid code"))
		c.Close()
		return
	}

	ProgressConnections[c] = true

	defer func() {
		delete(ProgressConnections, c)
		c.Close()
	}()

	for {
		// Keep the connection open.
		if _, _, err := c.ReadMessage(); err != nil {
			// When the ReadMessage method returns an error (e.g., when the client closes their browser),
			// the loop breaks, and the function returns.
			break
		}
	}
}

// BroadcastProgress sends a message to all WebSocket connections.
func BroadcastProgress(data responses.FileProgress) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	for c := range ProgressConnections {
		c.WriteMessage(websocket.TextMessage, message)
	}
}

// Handshake is a WebSocket handler that creates a unique code for the handshake.
func Handshake(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id := c.Query("id")
	app := c.Query("app")

	// Check if app exists.
	if available, err := services.IsAppAvailable(app); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Create unique code for the handshake.
	code, err := services.CreateHandshake(app, id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.CacheError, err.Error())
	}

	// send response.
	response := responses.Handshake{}
	response.SetHandshake(code)

	return c.JSON(response)
}
