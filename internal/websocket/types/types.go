package types

// WebSocket message types
const (
	// Card events
	MSG_CARD_CREATED = "card_created"
	MSG_CARD_UPDATED = "card_updated"
	MSG_CARD_MOVED   = "card_moved"
	MSG_CARD_DELETED = "card_deleted"

	// List events
	MSG_LIST_CREATED = "list_created"
	MSG_LIST_UPDATED = "list_updated"
	MSG_LIST_DELETED = "list_deleted"
	MSG_LIST_MOVED   = "list_moved"

	// User events
	MSG_USER_JOINED = "user_joined"
	MSG_USER_LEFT   = "user_left"
	MSG_USER_TYPING = "user_typing"

	// Board events
	MSG_BOARD_UPDATED = "board_updated"

	// System events
	MSG_ERROR = "error"
	MSG_PING  = "ping"
	MSG_PONG  = "pong"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string      `json:"type"`
	BoardID   string      `json:"board_id"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	UserID    string      `json:"user_id,omitempty"`
}

// BroadcastMessage represents a message to be broadcasted
type BroadcastMessage struct {
	BoardID string    `json:"board_id"`
	Message WSMessage `json:"message"`
}

// ClientMessage represents a message from client
type ClientMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
