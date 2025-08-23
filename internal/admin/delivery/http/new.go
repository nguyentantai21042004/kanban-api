package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type Handler interface {
	Dashboard(c *gin.Context)
	Users(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	Health(c *gin.Context)
	Roles(c *gin.Context)
}

type handler struct {
	l  pkgLog.Logger
	uc admin.UseCase
	d  *discord.Discord
}

func New(l pkgLog.Logger, uc admin.UseCase, d *discord.Discord) Handler {
	return handler{l: l, uc: uc, d: d}
}
