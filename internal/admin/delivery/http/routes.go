package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/middleware"
)

func MapAdminRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.GET("/dashboard", h.Dashboard)
	r.GET("/users", h.Users)
	r.POST("/users", h.CreateUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.GET("/roles", h.Roles)
	r.GET("/health", h.Health)
}
