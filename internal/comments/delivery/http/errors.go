package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/comments"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery = pkgErrors.NewHTTPError(10401, "Wrong query")
	// errWrongBody     = pkgErrors.NewHTTPError(10402, "Wrong body")
	errNotFound              = pkgErrors.NewHTTPError(10403, "Comment not found")
	errFieldRequired         = pkgErrors.NewHTTPError(10404, "Field required")
	errCardNotFound          = pkgErrors.NewHTTPError(10405, "Card not found")
	errParentCommentNotFound = pkgErrors.NewHTTPError(10406, "Parent comment not found")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case comments.ErrCommentNotFound:
		return errNotFound
	case comments.ErrFieldRequired:
		return errFieldRequired
	case comments.ErrCardNotFound:
		return errCardNotFound
	case comments.ErrParentCommentNotFound:
		return errParentCommentNotFound
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
	errCardNotFound,
	errParentCommentNotFound,
}
