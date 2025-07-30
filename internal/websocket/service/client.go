package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/types"
)

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	boardID string
	userID  string

	// Client metadata
	joinedAt time.Time
	lastSeen time.Time
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
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle client messages (ping, typing indicators, etc.)
		c.handleClientMessage(message)
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
		}
	}
}

// StartPumps starts the read and write pumps for the client
func (c *Client) StartPumps() {
	go c.writePump()
	go c.readPump()
}

// handleClientMessage processes incoming messages from the client
func (c *Client) handleClientMessage(message []byte) {
	var clientMsg types.ClientMessage

	if err := json.Unmarshal(message, &clientMsg); err != nil {
		log.Printf("Error parsing client message: %v", err)
		return
	}

	log.Printf("Received WebSocket message from client %s: type=%s, data=%+v", c.userID, clientMsg.Type, clientMsg.Data)

	switch clientMsg.Type {
	case types.MSG_AUTH:
		// Handle authentication message
		log.Printf("Client %s authenticated successfully", c.userID)
		// Send confirmation back to client
		authConfirmMsg := types.WSMessage{
			Type:      "auth_confirmed",
			BoardID:   c.boardID,
			Data:      map[string]interface{}{"status": "authenticated"},
			Timestamp: time.Now().Unix(),
			UserID:    c.userID,
		}
		authConfirmBytes, _ := json.Marshal(authConfirmMsg)
		select {
		case c.send <- authConfirmBytes:
		default:
		}

	case types.MSG_PING:
		// Respond with pong
		pongMsg := types.WSMessage{
			Type:      types.MSG_PONG,
			BoardID:   c.boardID,
			Data:      map[string]interface{}{"timestamp": time.Now().Unix()},
			Timestamp: time.Now().Unix(),
			UserID:    c.userID,
		}

		pongBytes, _ := json.Marshal(pongMsg)
		select {
		case c.send <- pongBytes:
		default:
		}

	case types.MSG_USER_TYPING:
		// Broadcast typing indicator to other users
		c.hub.BroadcastToBoard(c.boardID, types.MSG_USER_TYPING, map[string]interface{}{
			"user_id":   c.userID,
			"card_id":   clientMsg.Data["card_id"],
			"is_typing": clientMsg.Data["is_typing"],
		}, c.userID)

	case types.MSG_CARD_CREATED, types.MSG_CARD_UPDATED, types.MSG_CARD_MOVED, types.MSG_CARD_DELETED:
		// Broadcast card events to all users in board
		log.Printf("Broadcasting card event: %s to board %s", clientMsg.Type, c.boardID)
		c.hub.BroadcastToBoard(c.boardID, clientMsg.Type, clientMsg.Data, c.userID)

	case types.MSG_LIST_CREATED, types.MSG_LIST_UPDATED, types.MSG_LIST_MOVED, types.MSG_LIST_DELETED:
		// Broadcast list events to all users in board
		log.Printf("Broadcasting list event: %s to board %s", clientMsg.Type, c.boardID)
		c.hub.BroadcastToBoard(c.boardID, clientMsg.Type, clientMsg.Data, c.userID)

	case types.MSG_BOARD_UPDATED:
		// Broadcast board events to all users in board
		log.Printf("Broadcasting board event: %s to board %s", clientMsg.Type, c.boardID)
		c.hub.BroadcastToBoard(c.boardID, clientMsg.Type, clientMsg.Data, c.userID)

	default:
		log.Printf("Unknown message type: %s", clientMsg.Type)
	}
}
