package service

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

// Global hub instance
var WSHub *Hub

// InitWebSocketHub initializes the global WebSocket hub
func InitWebSocketHub(logger log.Logger) error {
	WSHub = NewHub(logger)

	ctx := context.Background()
	if err := WSHub.Start(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "WebSocket Hub initialized")
	return nil
}

// GetHub returns the global hub instance
func GetHub() *Hub {
	return WSHub
}

// BroadcastToBoard is a convenience function to broadcast messages
func BroadcastToBoard(boardID string, msgType string, data interface{}, userID string) error {
	if WSHub == nil {
		return nil
	}

	ctx := context.Background()
	return WSHub.BroadcastToBoard(ctx, boardID, msgType, data, userID)
}

// GetActiveUsersCount is a convenience function to get active users count
func GetActiveUsersCount(boardID string) (int, error) {
	if WSHub == nil {
		return 0, nil
	}

	ctx := context.Background()
	return WSHub.GetActiveUsersCount(ctx, boardID)
}
