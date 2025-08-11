package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/websocket"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

// Hub implements the websocket.Hub interface
type Hub struct {
	// Board ID -> Client connections
	boards map[string]map[*Client]bool

	// Channel operations
	register   chan *Client
	unregister chan *Client
	broadcast  chan websocket.BroadcastMessage

	// Logger
	logger log.Logger

	// Mutex for concurrent access
	mutex sync.RWMutex

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// NewHub creates a new WebSocket hub
func NewHub(logger log.Logger) *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		boards:     make(map[string]map[*Client]bool),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		broadcast:  make(chan websocket.BroadcastMessage, 256),
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		done:       make(chan struct{}),
	}
}

// Start starts the hub's main loop
func (h *Hub) Start(ctx context.Context) error {
	if h.ctx.Err() != nil {
		return websocket.ErrHubNotInitialized{}
	}

	go h.run()
	h.logger.Info(ctx, "WebSocket Hub started")
	return nil
}

// Stop stops the hub gracefully
func (h *Hub) Stop(ctx context.Context) error {
	if h.ctx.Err() != nil {
		return websocket.ErrHubNotInitialized{}
	}

	h.cancel()
	close(h.done)
	h.logger.Info(ctx, "WebSocket Hub stopped")
	return nil
}

// run starts the hub's main loop
func (h *Hub) run() {
	// Cleanup ticker for inactive connections
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			if err := h.registerClient(context.Background(), client); err != nil {
				h.logger.Error(context.Background(), "Failed to register client", "error", err)
			}

		case client := <-h.unregister:
			if err := h.unregisterClient(context.Background(), client); err != nil {
				h.logger.Error(context.Background(), "Failed to unregister client", "error", err)
			}

		case message := <-h.broadcast:
			if err := h.broadcastToBoard(context.Background(), message.BoardID, message.Message); err != nil {
				h.logger.Error(context.Background(), "Failed to broadcast message", "error", err, "board_id", message.BoardID)
			}

		case <-ticker.C:
			if err := h.cleanupInactiveClients(context.Background()); err != nil {
				h.logger.Error(context.Background(), "Failed to cleanup inactive clients", "error", err)
			}

		case <-h.ctx.Done():
			return
		}
	}
}

// RegisterClient registers a new client with the hub
func (h *Hub) RegisterClient(ctx context.Context, client websocket.Client) error {
	select {
	case h.register <- client.(*Client):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return websocket.ErrBroadcastFailed{}
	}
}

// UnregisterClient removes a client from the hub
func (h *Hub) UnregisterClient(ctx context.Context, client websocket.Client) error {
	select {
	case h.unregister <- client.(*Client):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return websocket.ErrBroadcastFailed{}
	}
}

// registerClient adds a new client to the hub (internal use)
func (h *Hub) registerClient(ctx context.Context, client *Client) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.boards[client.boardID] == nil {
		h.boards[client.boardID] = make(map[*Client]bool)
	}

	h.boards[client.boardID][client] = true

	// Get current active users count
	activeUsers := len(h.boards[client.boardID])

	h.logger.Info(ctx, "Client joined board", "user_id", client.userID, "board_id", client.boardID, "active_users", activeUsers)

	// Send welcome message to new client
	welcomeData := map[string]interface{}{
		"user_id":      client.userID,
		"joined_at":    client.joinedAt,
		"active_users": activeUsers,
	}

	// Broadcast user joined to all clients in board
	go func() {
		select {
		case h.broadcast <- websocket.BroadcastMessage{
			BoardID: client.boardID,
			Message: websocket.WSMessage{
				Type:      websocket.MSG_USER_JOINED,
				BoardID:   client.boardID,
				Data:      welcomeData,
				Timestamp: time.Now().Unix(),
				UserID:    client.userID,
			},
		}:
		case <-h.ctx.Done():
		}
	}()

	return nil
}

// unregisterClient removes a client from the hub (internal use)
func (h *Hub) unregisterClient(ctx context.Context, client *Client) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if clients, ok := h.boards[client.boardID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)

			// Clean up empty board
			if len(clients) == 0 {
				delete(h.boards, client.boardID)
			}

			activeUsers := len(clients)

			h.logger.Info(ctx, "Client left board", "user_id", client.userID, "board_id", client.boardID, "remaining_users", activeUsers)

			// Broadcast user left
			leaveData := map[string]interface{}{
				"user_id":      client.userID,
				"left_at":      time.Now(),
				"active_users": activeUsers,
			}

			go func() {
				select {
				case h.broadcast <- websocket.BroadcastMessage{
					BoardID: client.boardID,
					Message: websocket.WSMessage{
						Type:      websocket.MSG_USER_LEFT,
						BoardID:   client.boardID,
						Data:      leaveData,
						Timestamp: time.Now().Unix(),
						UserID:    client.userID,
					},
				}:
				case <-h.ctx.Done():
				}
			}()
		}
	}

	return nil
}

// BroadcastToBoard broadcasts a message to all clients in a board
func (h *Hub) BroadcastToBoard(ctx context.Context, boardID, msgType string, data interface{}, userID string) error {
	message := websocket.WSMessage{
		Type:      msgType,
		BoardID:   boardID,
		Data:      data,
		Timestamp: time.Now().Unix(),
		UserID:    userID,
	}

	select {
	case h.broadcast <- websocket.BroadcastMessage{BoardID: boardID, Message: message}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return websocket.ErrBroadcastFailed{}
	}
}

// broadcastToBoard sends a message to all clients in a board (internal use)
func (h *Hub) broadcastToBoard(ctx context.Context, boardID string, message websocket.WSMessage) error {
	h.mutex.RLock()
	clients := h.boards[boardID]
	h.mutex.RUnlock()

	if clients == nil {
		return websocket.ErrBoardNotFound{}
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send to all clients in the board
	for client := range clients {
		select {
		case client.send <- messageBytes:
			client.lastSeen = time.Now()
		default:
			// Client send buffer is full, close connection
			h.logger.Warn(ctx, "Client send buffer full, closing connection", "user_id", client.userID)
			delete(clients, client)
			close(client.send)
		}
	}

	return nil
}

// cleanupInactiveClients removes inactive clients
func (h *Hub) cleanupInactiveClients(ctx context.Context) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)

	for boardID, clients := range h.boards {
		for client := range clients {
			if client.lastSeen.Before(cutoff) {
				h.logger.Info(ctx, "Removing inactive client", "user_id", client.userID, "board_id", boardID)
				delete(clients, client)
				close(client.send)
			}
		}

		// Clean up empty boards
		if len(clients) == 0 {
			delete(h.boards, boardID)
		}
	}

	return nil
}

// GetActiveUsersCount returns the number of active users in a board
func (h *Hub) GetActiveUsersCount(ctx context.Context, boardID string) (int, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, ok := h.boards[boardID]; ok {
		return len(clients), nil
	}
	return 0, websocket.ErrBoardNotFound{}
}

// GetConnectedBoards returns a list of all connected board IDs
func (h *Hub) GetConnectedBoards(ctx context.Context) ([]string, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	boards := make([]string, 0, len(h.boards))
	for boardID := range h.boards {
		boards = append(boards, boardID)
	}
	return boards, nil
}
