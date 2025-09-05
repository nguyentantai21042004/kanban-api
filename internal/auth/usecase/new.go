package usecase

import (
	"github.com/nguyentantai21042004/kanban-api/internal/auth"
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	"github.com/nguyentantai21042004/kanban-api/pkg/encrypter"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/scope"
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
