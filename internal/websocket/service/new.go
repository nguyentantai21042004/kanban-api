package service

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Global hub instance
var WSHub *Hub

// InitWebSocketHub initializes the global WebSocket hub
func InitWebSocketHub() {
	WSHub = NewHub()

	go WSHub.Run()

	log.Println("WebSocket Hub initialized")
}

// GetHub returns the global hub instance
func GetHub() *Hub {
	return WSHub
}

// BroadcastToBoard is a convenience function to broadcast messages
func BroadcastToBoard(boardID string, msgType string, data interface{}, userID string) {
	if WSHub == nil {
		return
	}

	WSHub.BroadcastToBoard(boardID, msgType, data, userID)
}

// GetActiveUsersCount is a convenience function to get active users count
func GetActiveUsersCount(boardID string) int {
	if WSHub == nil {
		return 0
	}

	return WSHub.GetActiveUsersCount(boardID)
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, boardID, userID string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		boardID:  boardID,
		userID:   userID,
		joinedAt: time.Now(),
		lastSeen: time.Now(),
	}
}
