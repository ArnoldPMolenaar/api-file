package controllers

import (
	"api-file/main/src/dto/responses"
	"api-file/main/src/errors"
	"api-file/main/src/services"
	"encoding/json"
	"log"
	"strconv"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Handshake failed", "message": "`+err.Error()+`", "code": "`+errors.CodeExists+`"}`))
		c.Close()
		return
	} else if value != id {
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Handshake failed", "message": "invalid code", "code": "`+errors.CodeInvalid+`"}`))
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
func BroadcastProgress(data *responses.FileProgress) {
	message, err := json.Marshal(&data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	for c := range ProgressConnections {
		_ = c.WriteMessage(websocket.TextMessage, message)
	}
}

// Handshake is a WebSocket handler that creates a unique code for the handshake.
func Handshake(c *fiber.Ctx) error {
	// Get the ID from the URL.
	appStoragePathIdParam := c.Query("id")
	app := c.Query("app")
	appStoragePathId, err := strconv.ParseUint(appStoragePathIdParam, 0, 32)

	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "AppStoragePathId is invalid.")
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(app); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Check if storage path exists within the app.
	if id, err := services.GetStoragePathIDByApp(app); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if id == nil || *id != uint(appStoragePathId) {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.StoragePathExists, "Storage path does not exist within the app.")
	}

	// Create unique code for the handshake.
	code, err := services.CreateHandshake(app, uint(appStoragePathId))
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.CacheError, err.Error())
	}

	// send response.
	response := responses.Handshake{}
	response.SetHandshake(code)

	return c.JSON(response)
}
