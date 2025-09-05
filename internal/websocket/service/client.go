package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	wsPkg "github.com/nguyentantai21042004/kanban-api/internal/websocket"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
)

// Client implements the wsPkg.Client interface
type Client struct {
	hub     *Hub
	l       log.Logger
	conn    *websocket.Conn
	send    chan []byte
	boardID string
	userID  string
	id      string

	// Client metadata
	joinedAt time.Time
	lastSeen time.Time

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
	mutex  sync.RWMutex
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, logger log.Logger, conn *websocket.Conn, boardID, userID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		hub:      hub,
		l:        logger,
		conn:     conn,
		send:     make(chan []byte, 256),
		boardID:  boardID,
		userID:   userID,
		id:       generateClientID(),
		joinedAt: time.Now(),
		lastSeen: time.Now(),
		ctx:      ctx,
		cancel:   cancel,
		done:     make(chan struct{}),
	}
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// Connect establishes the WebSocket connection
func (c *Client) Connect(ctx context.Context) error {
	if c.conn == nil {
		return wsPkg.ErrConnectionFailed{}
	}

	c.lastSeen = time.Now()
	c.l.Info(ctx, "Client connected", "user_id", c.userID, "board_id", c.boardID)
	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	c.l.Info(ctx, "Client disconnected", "user_id", c.userID, "board_id", c.boardID)
	return nil
}

// IsConnected checks if the client is currently connected
func (c *Client) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.conn != nil
}

// SendMessage sends a message to the client
func (c *Client) SendMessage(ctx context.Context, message wsPkg.WSMessage) error {
	if !c.IsConnected() {
		return wsPkg.ErrConnectionFailed{}
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	select {
	case c.send <- messageBytes:
		c.lastSeen = time.Now()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return wsPkg.ErrBroadcastFailed{}
	}
}

// ReadMessage reads a message from the client
func (c *Client) ReadMessage(ctx context.Context) (wsPkg.ClientMessage, error) {
	if !c.IsConnected() {
		return wsPkg.ClientMessage{}, wsPkg.ErrConnectionFailed{}
	}

	_, messageBytes, err := c.conn.ReadMessage()
	if err != nil {
		return wsPkg.ClientMessage{}, err
	}

	var clientMsg wsPkg.ClientMessage
	if err := json.Unmarshal(messageBytes, &clientMsg); err != nil {
		return wsPkg.ClientMessage{}, err
	}

	c.lastSeen = time.Now()
	return clientMsg, nil
}

// GetID returns the client's unique ID
func (c *Client) GetID() string {
	return c.id
}

// GetBoardID returns the board ID the client is connected to
func (c *Client) GetBoardID() string {
	return c.boardID
}

// GetUserID returns the user ID of the client
func (c *Client) GetUserID() string {
	return c.userID
}

// GetLastSeen returns the timestamp when the client was last seen
func (c *Client) GetLastSeen() int64 {
	return c.lastSeen.Unix()
}

// Start starts the client's read and write pumps
func (c *Client) Start(ctx context.Context) error {
	if c.ctx.Err() != nil {
		return wsPkg.ErrConnectionFailed{}
	}

	go c.writePump()
	go c.readPump()

	c.l.Info(ctx, "Client started", "user_id", c.userID, "board_id", c.boardID)
	return nil
}

// Stop stops the client gracefully
func (c *Client) Stop(ctx context.Context) error {
	if c.ctx.Err() != nil {
		return wsPkg.ErrConnectionFailed{}
	}

	c.cancel()
	close(c.done)

	if err := c.Disconnect(ctx); err != nil {
		return err
	}

	c.l.Info(ctx, "Client stopped", "user_id", c.userID, "board_id", c.boardID)
	return nil
}

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	const (
		pongWait       = 60 * time.Second
		pingPeriod     = (pongWait * 9) / 10
		maxMessageSize = 512
	)

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.lastSeen = time.Now()
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.l.Error(context.Background(), "WebSocket error", "error", err)
				}
				// On any read error, exit the read pump to avoid repeated reads on a failed connection
				return
			}

			// Handle client messages (ping, typing indicators, etc.)
			if err := c.handleClientMessage(message); err != nil {
				c.l.Error(context.Background(), "Failed to handle client message", "error", err)
			}
		}
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {
	const (
		writeWait  = 10 * time.Second
		pingPeriod = 54 * time.Second
	)

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

// handleClientMessage processes incoming messages from the client
func (c *Client) handleClientMessage(message []byte) error {
	var clientMsg wsPkg.ClientMessage

	if err := json.Unmarshal(message, &clientMsg); err != nil {
		return err
	}

	c.l.Info(context.Background(), "Received WebSocket message from client", "user_id", c.userID, "type", clientMsg.Type, "data", clientMsg.Data)

	switch clientMsg.Type {
	case wsPkg.MSG_AUTH:
		// Handle authentication message
		c.l.Info(context.Background(), "Client authenticated successfully", "user_id", c.userID)
		// Send confirmation back to client
		authConfirmMsg := wsPkg.WSMessage{
			Type:      "auth_confirmed",
			BoardID:   c.boardID,
			Data:      map[string]interface{}{"status": "authenticated"},
			Timestamp: time.Now().Unix(),
			UserID:    c.userID,
		}
		authConfirmBytes, _ := json.Marshal(authConfirmMsg)
		select {
		case c.send <- authConfirmBytes:
		case <-c.ctx.Done():
		}

	case wsPkg.MSG_PING:
		// Respond with pong
		pongMsg := wsPkg.WSMessage{
			Type:      wsPkg.MSG_PONG,
			BoardID:   c.boardID,
			Data:      map[string]interface{}{"timestamp": time.Now().Unix()},
			Timestamp: time.Now().Unix(),
			UserID:    c.userID,
		}

		pongBytes, _ := json.Marshal(pongMsg)
		select {
		case c.send <- pongBytes:
		case <-c.ctx.Done():
		}

	case wsPkg.MSG_USER_TYPING:
		// Broadcast typing indicator to other users
		if err := c.hub.BroadcastToBoard(context.Background(), c.boardID, wsPkg.MSG_USER_TYPING, map[string]interface{}{
			"user_id":   c.userID,
			"card_id":   clientMsg.Data["card_id"],
			"is_typing": clientMsg.Data["is_typing"],
		}, c.userID); err != nil {
			c.l.Error(context.Background(), "Failed to broadcast typing indicator", "error", err)
		}

	case wsPkg.MSG_CARD_CREATED, wsPkg.MSG_CARD_UPDATED, wsPkg.MSG_CARD_MOVED, wsPkg.MSG_CARD_DELETED:
		// Broadcast card events to all users in board
		c.l.Info(context.Background(), "Broadcasting card event", "type", clientMsg.Type, "board_id", c.boardID)
		if err := c.hub.BroadcastToBoard(context.Background(), c.boardID, clientMsg.Type, clientMsg.Data, c.userID); err != nil {
			c.l.Error(context.Background(), "Failed to broadcast card event", "error", err)
		}

	case wsPkg.MSG_LIST_CREATED, wsPkg.MSG_LIST_UPDATED, wsPkg.MSG_LIST_MOVED, wsPkg.MSG_LIST_DELETED:
		// Broadcast list events to all users in board
		c.l.Info(context.Background(), "Broadcasting list event", "type", clientMsg.Type, "board_id", c.boardID)
		if err := c.hub.BroadcastToBoard(context.Background(), c.boardID, clientMsg.Type, clientMsg.Data, c.userID); err != nil {
			c.l.Error(context.Background(), "Failed to broadcast list event", "error", err)
		}

	case wsPkg.MSG_BOARD_UPDATED:
		// Broadcast board events to all users in board
		c.l.Info(context.Background(), "Broadcasting board event", "type", clientMsg.Type, "board_id", c.boardID)
		if err := c.hub.BroadcastToBoard(context.Background(), c.boardID, clientMsg.Type, clientMsg.Data, c.userID); err != nil {
			c.l.Error(context.Background(), "Failed to broadcast board event", "error", err)
		}

	default:
		c.l.Warn(context.Background(), "Unknown message type", "type", clientMsg.Type)
	}

	return nil
}
