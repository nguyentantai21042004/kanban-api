package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/user"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery    = pkgErrors.NewHTTPError(10601, "Wrong query")
	errNotFound      = pkgErrors.NewHTTPError(10603, "User not found")
	errFieldRequired = pkgErrors.NewHTTPError(10604, "Field required")
	errUnauthorized  = pkgErrors.NewHTTPError(10605, "Unauthorized")
	errUserExists    = pkgErrors.NewHTTPError(10606, "User already exists")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case user.ErrUserNotFound:
		return errNotFound
	case user.ErrUserExists:
		return errUserExists
	case user.ErrUnauthorized:
		return errUnauthorized
	case user.ErrFieldRequired:
		return errFieldRequired
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
}
