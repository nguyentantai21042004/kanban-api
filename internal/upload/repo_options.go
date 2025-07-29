package upload

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/models"
	pag "gitlab.com/tantai-kanban/kanban-api/pkg/paginator"
)

type CreateOptions struct {
	Upload models.Upload
}

type GetOptions struct {
	Filter   Filter
	PagQuery pag.PaginateQuery
}
