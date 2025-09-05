package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.POST("/login", h.Login)
	r.POST("/refresh", h.RefreshToken)
	r.POST("/logout", mw.Auth(), h.Logout)
}
