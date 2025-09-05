package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/labels"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery = pkgErrors.NewHTTPError(10201, "Wrong query")
	// errWrongBody     = pkgErrors.NewHTTPError(10202, "Wrong body")
	errNotFound      = pkgErrors.NewHTTPError(10203, "Label not found")
	errFieldRequired = pkgErrors.NewHTTPError(10204, "Field required")
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
