package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/lists"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery    = pkgErrors.NewHTTPError(10101, "Wrong query")
	errWrongBody     = pkgErrors.NewHTTPError(10102, "Wrong body")
	errNotFound      = pkgErrors.NewHTTPError(10103, "List not found")
	errFieldRequired = pkgErrors.NewHTTPError(10104, "Field required")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case lists.ErrNotFound:
		return errNotFound
	case lists.ErrFieldRequired:
		return errFieldRequired
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
}
