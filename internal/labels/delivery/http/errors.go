package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/labels"
	pkgErrors "gitlab.com/tantai-kanban/kanban-api/pkg/errors"
)

var (
	errWrongQuery = pkgErrors.NewHTTPError(10001, "Wrong query")
	// errWrongBody     = pkgErrors.NewHTTPError(10002, "Wrong body")
	errNotFound      = pkgErrors.NewHTTPError(10003, "Label not found")
	errFieldRequired = pkgErrors.NewHTTPError(10004, "Field required")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case labels.ErrLabelNotFound:
		return errNotFound
	case labels.ErrFieldRequired:
		return errFieldRequired
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
}
