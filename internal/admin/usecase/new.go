package usecase

import (
	"github.com/nguyentantai21042004/kanban-api/internal/admin"
	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	"github.com/nguyentantai21042004/kanban-api/internal/comments"
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/internal/websocket"
	pkgLog "github.com/nguyentantai21042004/kanban-api/pkg/log"
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
