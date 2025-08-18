package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/admin"
	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
	"gitlab.com/tantai-kanban/kanban-api/internal/comments"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/internal/websocket"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type implUsecase struct {
	l         pkgLog.Logger
	userUC    user.UseCase
	boardUC   boards.UseCase
	cardUC    cards.UseCase
	commentUC comments.UseCase
	roleUC    role.UseCase
	wsHub     websocket.Hub
}

func New(l pkgLog.Logger, userUC user.UseCase, boardUC boards.UseCase, cardUC cards.UseCase, commentUC comments.UseCase, roleUC role.UseCase, wsHub websocket.Hub) admin.UseCase {
	return implUsecase{l: l, userUC: userUC, boardUC: boardUC, cardUC: cardUC, commentUC: commentUC, roleUC: roleUC, wsHub: wsHub}
}
