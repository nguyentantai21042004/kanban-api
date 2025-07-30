package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

type Handler struct {
	hub        *service.Hub
	jwtManager scope.Manager
}

// New creates a new WebSocket handler
func New(hub *service.Hub, jwtManager scope.Manager) *Handler {
	return &Handler{
		hub:        hub,
		jwtManager: jwtManager,
	}
}

// ServeWebSocket handles WebSocket upgrade and connection
// @Summary WebSocket Connection
// @Description Establish WebSocket connection for real-time collaboration on a board
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param board_id path string true "Board ID" example("board123")
// @Param token query string true "JWT token" example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
// @Success 101 "Switching Protocols" {string} string "WebSocket connection established"
// @Failure 400 {object} map[string]interface{} "Bad Request - board_id is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized - user authentication required"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/websocket/ws/{board_id} [get]
func (h *Handler) ServeWebSocket(c *gin.Context) {
	boardID := c.Param("board_id")
	token := c.Query("token")

	log.Printf("WebSocket connection attempt - BoardID: %s, Token: %s", boardID, token[:50]+"...")

	if boardID == "" {
		log.Printf("WebSocket error: board_id is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}

	if token == "" {
		log.Printf("WebSocket error: token is required")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	// Validate JWT token
	log.Printf("Validating JWT token...")
	payload, err := h.jwtManager.Verify(token)
	if err != nil {
		log.Printf("WebSocket JWT validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	log.Printf("JWT validation successful - UserID: %s", payload.UserID)

	userID := payload.UserID
	if userID == "" {
		log.Printf("WebSocket error: user_id not found in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	log.Printf("WebSocket connection authorized for user %s on board %s", userID, boardID)

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
