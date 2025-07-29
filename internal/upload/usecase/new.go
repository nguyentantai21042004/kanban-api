package usecase

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/upload"
	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/minio"
)

type usecase struct {
	l     log.Logger
	repo  upload.Repository
	minio minio.MinIO
}

func New(l log.Logger, repo upload.Repository, minio minio.MinIO) upload.UseCase {
	return &usecase{
		l:     l,
		repo:  repo,
		minio: minio,
	}
}
