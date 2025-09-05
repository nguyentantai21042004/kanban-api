package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/pkg/discord"
	pkgLog "github.com/nguyentantai21042004/kanban-api/pkg/log"
)

type Handler interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Detail(c *gin.Context)
	Delete(c *gin.Context)
	Move(c *gin.Context)
	GetActivities(c *gin.Context)

	// Enhanced functionality methods
	Assign(c *gin.Context)
	Unassign(c *gin.Context)
	AddAttachment(c *gin.Context)
	RemoveAttachment(c *gin.Context)
	UpdateTimeTracking(c *gin.Context)
	UpdateChecklist(c *gin.Context)
	AddTag(c *gin.Context)
	RemoveTag(c *gin.Context)
	SetStartDate(c *gin.Context)
	SetCompletionDate(c *gin.Context)
}

type handler struct {
	l  pkgLog.Logger
	uc cards.UseCase
	d  *discord.Discord
}

func New(l pkgLog.Logger, uc cards.UseCase, d *discord.Discord) Handler {
	h := handler{
		l:  l,
		uc: uc,
		d:  d,
	}
	return h
}
