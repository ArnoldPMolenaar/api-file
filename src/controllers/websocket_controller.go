package controllers

import "github.com/gofiber/contrib/websocket"

// ProgressConnections is a map of WebSocket connections from clients.
var ProgressConnections = make(map[*websocket.Conn]bool)

// WebSocketProgress is a WebSocket handler that sends progress updates to the client.
func WebSocketProgress(c *websocket.Conn) {
	ProgressConnections[c] = true

	defer func() {
		delete(ProgressConnections, c)
		c.Close()
	}()

	for {
		// Keep the connection open
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
