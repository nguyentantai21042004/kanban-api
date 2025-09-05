package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/boards"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery = pkgErrors.NewHTTPError(10301, "Wrong query")
	// errWrongBody     = pkgErrors.NewHTTPError(10302, "Wrong body")
	errNotFound      = pkgErrors.NewHTTPError(10303, "Board not found")
	errFieldRequired = pkgErrors.NewHTTPError(10304, "Field required")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case boards.ErrNotFound:
		return errNotFound
	case boards.ErrFieldRequired:
		return errFieldRequired
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
}
