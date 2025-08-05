package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
)

func MapCommentRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.GET("", h.Get)
	r.POST("", h.Create)
	r.PUT("/:id", h.Update)
	r.GET("/:id", h.Detail)
	r.DELETE("", h.Delete)
}

func MapCardCommentRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.GET("/:card_id/comments", h.GetByCard)
}
