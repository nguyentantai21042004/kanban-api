package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/auth"
	"gitlab.com/tantai-kanban/kanban-api/internal/user"
	"gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type usecase struct {
	l       log.Logger
	userUC  user.UseCase
	encrypt encrypter.Encrypter
}

func New(l log.Logger, userUC user.UseCase, encrypt encrypter.Encrypter) auth.UseCase {
	return &usecase{
		l:       l,
		userUC:  userUC,
		encrypt: encrypt,
	}
}
