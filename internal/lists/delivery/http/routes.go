package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
)

func MapListRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.GET("", h.Get)
	r.POST("", h.Create)
	r.PUT("", h.Update)
	r.GET("/:id", h.Detail)
	r.DELETE("", h.Delete)
	// Move list
	r.POST("/move", h.Move)
}
