package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/middleware"
)

func MapCardRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.GET("", h.Get)
	r.POST("", h.Create)
	r.PUT("", h.Update)
	r.GET("/:id", h.Detail)
	r.DELETE("", h.Delete)
	r.POST("/move", h.Move)
	r.GET("/activities", h.GetActivities)

	// Enhanced functionality routes
	r.POST("/assign", h.Assign)
	r.POST("/unassign", h.Unassign)
	r.POST("/attachments/add", h.AddAttachment)
	r.POST("/attachments/remove", h.RemoveAttachment)
	r.PUT("/time-tracking", h.UpdateTimeTracking)
	r.PUT("/checklist", h.UpdateChecklist)
	r.POST("/tags/add", h.AddTag)
	r.POST("/tags/remove", h.RemoveTag)
	r.PUT("/start-date", h.SetStartDate)
	r.PUT("/completion-date", h.SetCompletionDate)
}
