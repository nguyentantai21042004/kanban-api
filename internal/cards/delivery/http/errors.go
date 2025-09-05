package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/cards"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery    = pkgErrors.NewHTTPError(10001, "Wrong query")
	errWrongBody     = pkgErrors.NewHTTPError(10002, "Wrong body")
	errNotFound      = pkgErrors.NewHTTPError(10003, "Card not found")
	errFieldRequired = pkgErrors.NewHTTPError(10004, "Field required")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case cards.ErrCardNotFound:
		return errNotFound
	case cards.ErrFieldRequired:
		return errFieldRequired
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
}
