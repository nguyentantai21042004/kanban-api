package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/comments"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type Handler interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Detail(c *gin.Context)
	Delete(c *gin.Context)
	GetByCard(c *gin.Context)
}

type handler struct {
	l  pkgLog.Logger
	uc comments.UseCase
	d  *discord.Discord
}

func New(l pkgLog.Logger, uc comments.UseCase, d *discord.Discord) Handler {
	h := handler{
		l:  l,
		uc: uc,
		d:  d,
	}
	return h
}
