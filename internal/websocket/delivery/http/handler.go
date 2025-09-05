package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	wsPkg "github.com/nguyentantai21042004/kanban-api/internal/websocket"
	wsService "github.com/nguyentantai21042004/kanban-api/internal/websocket/service"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/scope"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

// Handler implements the wsPkg.Handler interface
type Handler struct {
	hub        wsPkg.Hub
	jwtManager scope.Manager
	logger     log.Logger
}

// New creates a new WebSocket handler
func New(hub wsPkg.Hub, jwtManager scope.Manager, logger log.Logger) *Handler {
	return &Handler{
		hub:        hub,
		jwtManager: jwtManager,
		logger:     logger,
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

	h.logger.Info(c.Request.Context(), "WebSocket connection attempt", "board_id", boardID, "token_prefix", token[:50]+"...")

	if boardID == "" {
		h.logger.Error(c.Request.Context(), "WebSocket error: board_id is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}

	if token == "" {
		h.logger.Error(c.Request.Context(), "WebSocket error: token is required")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	// Validate JWT token
	h.logger.Info(c.Request.Context(), "Validating JWT token")
	payload, err := h.jwtManager.Verify(token)
	if err != nil {
		h.logger.Error(c.Request.Context(), "WebSocket JWT validation failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	h.logger.Info(c.Request.Context(), "JWT validation successful", "user_id", payload.UserID)

	userID := payload.UserID
	if userID == "" {
		h.logger.Error(c.Request.Context(), "WebSocket error: user_id not found in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	h.logger.Info(c.Request.Context(), "WebSocket connection authorized", "user_id", userID, "board_id", boardID)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error(c.Request.Context(), "WebSocket upgrade failed", "error", err)
		return
	}

	// Create new client using constructor
	hub, ok := h.hub.(*wsService.Hub)
	if !ok {
		h.logger.Error(c.Request.Context(), "Invalid hub type")
		conn.Close()
		return
	}
	client := wsService.NewClient(hub, h.logger, conn, boardID, userID)

	// Register client with hub
	if err := h.hub.RegisterClient(c.Request.Context(), client); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to register client", "error", err)
		conn.Close()
		return
	}

	// Start client
	if err := client.Start(c.Request.Context()); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to start client", "error", err)
		conn.Close()
		return
	}
}
