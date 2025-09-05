package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/middleware"
)

// MapWebSocketRoutes maps WebSocket routes
func MapWebSocketRoutes(r *gin.RouterGroup, h *Handler, mw middleware.Middleware) {
	// WebSocket route requires authentication
	r.GET("/ws/:board_id", h.ServeWebSocket)
}
