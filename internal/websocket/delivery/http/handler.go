package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

type Handler struct {
	hub *service.Hub
}

// New creates a new WebSocket handler
func New(hub *service.Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// ServeWebSocket handles WebSocket upgrade and connection
// @Summary WebSocket Connection
// @Description Establish WebSocket connection for real-time collaboration on a board
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param board_id path string true "Board ID" example("board123")
// @Param Authorization header string true "Bearer token" example("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
// @Success 101 "Switching Protocols" {string} string "WebSocket connection established"
// @Failure 400 {object} map[string]interface{} "Bad Request - board_id is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized - user authentication required"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/websocket/ws/{board_id} [get]
func (h *Handler) ServeWebSocket(c *gin.Context) {
	boardID := c.Param("board_id")
	userID := c.GetString("user_id") // From JWT middleware

	if boardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user authentication required"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create new client using constructor
	client := service.NewClient(h.hub, conn, boardID, userID)

	// Register client with hub
	h.hub.RegisterClient(client)

	// Start client pumps
	client.StartPumps()
}
