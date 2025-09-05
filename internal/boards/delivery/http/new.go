package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/pkg/discord"
	pkgLog "github.com/nguyentantai21042004/kanban-api/pkg/log"
)

type Handler interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Detail(c *gin.Context)
	Delete(c *gin.Context)
}

type handler struct {
	l  pkgLog.Logger
	uc boards.UseCase
	d  *discord.Discord
}

func New(l pkgLog.Logger, uc boards.UseCase, d *discord.Discord) Handler {
	h := handler{
		l:  l,
		uc: uc,
		d:  d,
	}
	return h
}
