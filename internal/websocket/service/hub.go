package service

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/types"
)

type Hub struct {
	// Board ID -> Client connections
	boards map[string]map[*Client]bool

	// Channel operations
	register   chan *Client
	unregister chan *Client
	broadcast  chan types.BroadcastMessage

	// Mutex for concurrent access
	mutex sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		boards:     make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan types.BroadcastMessage, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	// Cleanup ticker for inactive connections
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToBoard(message.BoardID, message.Message)

		case <-ticker.C:
			h.cleanupInactiveClients()
		}
	}
}

// registerClient adds a new client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.boards[client.boardID] == nil {
		h.boards[client.boardID] = make(map[*Client]bool)
	}

	h.boards[client.boardID][client] = true

	// Get current active users count
	activeUsers := len(h.boards[client.boardID])

	log.Printf("Client %s joined board %s. Active users: %d",
		client.userID, client.boardID, activeUsers)

	// Send welcome message to new client
	welcomeData := map[string]interface{}{
		"user_id":      client.userID,
		"joined_at":    client.joinedAt,
		"active_users": activeUsers,
	}

	// Broadcast user joined to all clients in board
	go func() {
		h.broadcast <- types.BroadcastMessage{
			BoardID: client.boardID,
			Message: types.WSMessage{
				Type:      types.MSG_USER_JOINED,
				BoardID:   client.boardID,
				Data:      welcomeData,
				Timestamp: time.Now().Unix(),
				UserID:    client.userID,
			},
		}
	}()
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
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

			log.Printf("Client %s left board %s. Remaining users: %d",
				client.userID, client.boardID, activeUsers)

			// Broadcast user left
			leaveData := map[string]interface{}{
				"user_id":      client.userID,
				"left_at":      time.Now(),
				"active_users": activeUsers,
			}

			go func() {
				h.broadcast <- types.BroadcastMessage{
					BoardID: client.boardID,
					Message: types.WSMessage{
						Type:      types.MSG_USER_LEFT,
						BoardID:   client.boardID,
						Data:      leaveData,
						Timestamp: time.Now().Unix(),
						UserID:    client.userID,
					},
				}
			}()
		}
	}
}

// broadcastToBoard sends a message to all clients in a board
func (h *Hub) broadcastToBoard(boardID string, message types.WSMessage) {
	h.mutex.RLock()
	clients := h.boards[boardID]
	h.mutex.RUnlock()

	if clients == nil {
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	// Send to all clients in the board
	for client := range clients {
		select {
		case client.send <- messageBytes:
			client.lastSeen = time.Now()
		default:
			// Client send buffer is full, close connection
			log.Printf("Client %s send buffer full, closing connection", client.userID)
			delete(clients, client)
			close(client.send)
		}
	}
}

// cleanupInactiveClients removes inactive clients
func (h *Hub) cleanupInactiveClients() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)

	for boardID, clients := range h.boards {
		for client := range clients {
			if client.lastSeen.Before(cutoff) {
				log.Printf("Removing inactive client %s from board %s", client.userID, boardID)
				delete(clients, client)
				close(client.send)
			}
		}

		// Clean up empty boards
		if len(clients) == 0 {
			delete(h.boards, boardID)
		}
	}
}

// BroadcastToBoard is a public function for broadcasting messages to a board
func (h *Hub) BroadcastToBoard(boardID string, msgType string, data interface{}, userID string) {
	message := types.WSMessage{
		Type:      msgType,
		BoardID:   boardID,
		Data:      data,
		Timestamp: time.Now().Unix(),
		UserID:    userID,
	}

	select {
	case h.broadcast <- types.BroadcastMessage{BoardID: boardID, Message: message}:
	default:
		log.Printf("Broadcast channel full, dropping message for board %s", boardID)
	}
}

// RegisterClient registers a new client with the hub
func (h *Hub) RegisterClient(client *Client) {
	select {
	case h.register <- client:
	default:
		// Channel is full, log error
		log.Printf("Register channel full, dropping client registration")
	}
}

// GetActiveUsersCount returns the number of active users in a board
func (h *Hub) GetActiveUsersCount(boardID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, ok := h.boards[boardID]; ok {
		return len(clients)
	}
	return 0
}
