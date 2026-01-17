package app

import (
	"log"
	"net/http"
	"strings"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin (adjust for production with proper CORS)
		return true
	},
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *util.Hub, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		// Get user ID from query parameter or Authorization header
		userID := c.Query("user_id")
		if userID == "" {
			// Try to get from Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token := parts[1]
					claims, err := util.ValidateToken(token, jwtSecret)
					if err == nil {
						userID = claims.UserID
					}
				}
			}
		}

		if userID == "" {
			log.Println("WebSocket connection rejected: no user ID")
			conn.WriteJSON(gin.H{
				"type":    "error",
				"message": "User ID required",
			})
			conn.Close()
			return
		}

		// Create new client
		client := &util.Client{
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			UserID: userID,
		}

		// Register client with hub
		hub.Register <- client

		// Start goroutines for reading and writing
		go client.WritePump()
		go client.ReadPump()

		log.Printf("WebSocket client connected for user %s", userID)
	}
}
