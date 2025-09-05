package service

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/websocket"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
)

// WebSocketService implements the websocket.Service interface
type WebSocketService struct {
	hub    *Hub
	logger log.Logger
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(logger log.Logger) *WebSocketService {
	return &WebSocketService{
		logger: logger,
	}
}

// InitHub initializes the WebSocket hub
func (s *WebSocketService) InitHub(ctx context.Context, logger log.Logger) error {
	s.hub = NewHub(logger)
	return s.hub.Start(ctx)
}

// GetHub returns the WebSocket hub
func (s *WebSocketService) GetHub() websocket.Hub {
	return s.hub
}

// Shutdown gracefully shuts down the WebSocket service
func (s *WebSocketService) Shutdown(ctx context.Context) error {
	if s.hub != nil {
		return s.hub.Stop(ctx)
	}
	return nil
}

// BroadcastToBoard broadcasts a message to all clients in a board
func (s *WebSocketService) BroadcastToBoard(ctx context.Context, boardID, msgType string, data interface{}, userID string) error {
	if s.hub == nil {
		return websocket.ErrHubNotInitialized{}
	}
	return s.hub.BroadcastToBoard(ctx, boardID, msgType, data, userID)
}

// GetActiveUsersCount returns the number of active users in a board
func (s *WebSocketService) GetActiveUsersCount(ctx context.Context, boardID string) (int, error) {
	if s.hub == nil {
		return 0, websocket.ErrHubNotInitialized{}
	}
	return s.hub.GetActiveUsersCount(ctx, boardID)
}
