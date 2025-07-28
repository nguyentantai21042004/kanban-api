package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
)

// MapWebSocketRoutes maps WebSocket routes
func MapWebSocketRoutes(r *gin.RouterGroup, h *Handler, mw middleware.Middleware) {
	// WebSocket route requires authentication
	r.Use(mw.Auth())
	r.GET("/ws/:boardID", h.ServeWebSocket)
}
