package usecase

import (
	"github.com/nguyentantai21042004/kanban-api/internal/upload"
	"github.com/nguyentantai21042004/kanban-api/pkg/log"
	"github.com/nguyentantai21042004/kanban-api/pkg/minio"
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
