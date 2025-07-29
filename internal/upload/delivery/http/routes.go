package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapUploadRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.POST("", mw.Auth(), h.Create)
	r.GET("", mw.Auth(), h.Get)
	r.GET("/:id", mw.Auth(), h.Detail)
}
