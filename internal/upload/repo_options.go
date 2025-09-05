package upload

import (
	"github.com/nguyentantai21042004/kanban-api/internal/models"
	pag "github.com/nguyentantai21042004/kanban-api/pkg/paginator"
)

type CreateOptions struct {
	Upload models.Upload
}

type GetOptions struct {
	Filter   Filter
	PagQuery pag.PaginateQuery
}
