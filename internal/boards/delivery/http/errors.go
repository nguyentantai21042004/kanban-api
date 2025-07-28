package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/boards"
	pkgErrors "gitlab.com/tantai-kanban/kanban-api/pkg/errors"
)

var (
	errWrongQuery = pkgErrors.NewHTTPError(10001, "Wrong query")
	// errWrongBody     = pkgErrors.NewHTTPError(10002, "Wrong body")
	errNotFound      = pkgErrors.NewHTTPError(10003, "Board not found")
	errFieldRequired = pkgErrors.NewHTTPError(10004, "Field required")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case boards.ErrBoardNotFound:
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
