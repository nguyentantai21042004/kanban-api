package upload

import (
	"mime/multipart"

	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	pag "gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type CreateInput struct {
	FileHeader *multipart.FileHeader
	From       string
	BucketName string
}

type GetInput struct {
	Filter   Filter
	PagQuery pag.PaginateQuery
}

type Filter struct {
	ID            *string `json:"id"`
	BucketName    *string `json:"bucket_name"`
	OriginalName  *string `json:"original_name"`
	Source        *string `json:"source"`
	CreatedUserID *string `json:"created_user_id"`
}

type UploadOutput struct {
	Upload models.Upload
}

type GetOutput struct {
	Uploads   []models.Upload
	Paginator pag.Paginator
}

var FromTypes = []string{
	MinIO,
}

const (
	MinIO = "minio"
)
