package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapUserRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.GET("/me", mw.Auth(), h.DetailMe)
	r.PUT("/profile", mw.Auth(), h.UpdateProfile)
	r.GET("/:id", mw.Auth(), h.Detail)
	r.POST("", mw.Auth(), h.Create)
}
