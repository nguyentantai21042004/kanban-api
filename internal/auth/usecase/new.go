package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/auth"
	"gitlab.com/tantai-kanban/kanban-api/internal/role"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
)

type implUseCase struct {
	l       log.Logger
	encrypt encrypter.Encrypter
	scopeUC scope.Manager
	userUC  user.UseCase
	roleUC  role.UseCase
}

var _ auth.UseCase = &implUseCase{}

func New(l log.Logger, encrypt encrypter.Encrypter, scopeUC scope.Manager, userUC user.UseCase, roleUC role.UseCase) auth.UseCase {
	return &implUseCase{
		l:       l,
		encrypt: encrypt,
		scopeUC: scopeUC,
		userUC:  userUC,
		roleUC:  roleUC,
	}
}
