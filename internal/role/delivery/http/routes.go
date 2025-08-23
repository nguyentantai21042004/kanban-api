package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapRoleRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.GET("/", mw.Auth(), h.List)
	r.GET("/:id", mw.Auth(), h.Detail)
}
