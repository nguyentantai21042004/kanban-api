package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/middleware"
)

func MapBoardRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.GET("", h.Get)
	r.POST("", h.Create)
	r.PUT("", h.Update)
	r.GET("/:id", h.Detail)
	r.DELETE("", h.Delete)
}
