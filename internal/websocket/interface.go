package websocket

import (
	"context"
	"net/http"
)

// Service defines the interface for WebSocket operations
type Service interface {
	// Hub management
	InitHub(ctx context.Context, logger Logger) error
	GetHub() Hub
	Shutdown(ctx context.Context) error

	// Broadcasting
	BroadcastToBoard(ctx context.Context, boardID, msgType string, data interface{}, userID string) error
	GetActiveUsersCount(ctx context.Context, boardID string) (int, error)
}

// Hub defines the interface for WebSocket hub operations
type Hub interface {
	// Client management
	RegisterClient(ctx context.Context, client Client) error
	UnregisterClient(ctx context.Context, client Client) error

	// Broadcasting
	BroadcastToBoard(ctx context.Context, boardID, msgType string, data interface{}, userID string) error

	// Information
	GetActiveUsersCount(ctx context.Context, boardID string) (int, error)
	GetConnectedBoards(ctx context.Context) ([]string, error)

	// Lifecycle
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Client defines the interface for WebSocket client operations
type Client interface {
	// Connection management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	IsConnected() bool

	// Message handling
	SendMessage(ctx context.Context, message WSMessage) error
	ReadMessage(ctx context.Context) (ClientMessage, error)

	// Information
	GetID() string
	GetBoardID() string
	GetUserID() string
	GetLastSeen() int64

	// Lifecycle
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Handler defines the interface for WebSocket HTTP handler
type Handler interface {
	ServeWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

// Logger defines the interface for logging operations
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Helper function to create fields
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Error types
type (
	ErrHubNotInitialized struct{}
	ErrClientNotFound    struct{}
	ErrBoardNotFound     struct{}
	ErrConnectionFailed  struct{}
	ErrMessageInvalid    struct{}
	ErrBroadcastFailed   struct{}
)

func (e ErrHubNotInitialized) Error() string { return "websocket hub not initialized" }
func (e ErrClientNotFound) Error() string    { return "websocket client not found" }
func (e ErrBoardNotFound) Error() string     { return "websocket board not found" }
func (e ErrConnectionFailed) Error() string  { return "websocket connection failed" }
func (e ErrMessageInvalid) Error() string    { return "websocket message invalid" }
func (e ErrBroadcastFailed) Error() string   { return "websocket broadcast failed" }
