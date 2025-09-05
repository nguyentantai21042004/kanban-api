package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/auth"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery         = pkgErrors.NewHTTPError(10701, "Wrong query")
	errInvalidCredentials = pkgErrors.NewHTTPError(10702, "Invalid credentials")
	errInvalidToken       = pkgErrors.NewHTTPError(10703, "Invalid token")
	errTokenExpired       = pkgErrors.NewHTTPError(10704, "Token expired")
	errUnauthorized       = pkgErrors.NewHTTPError(10705, "Unauthorized")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case auth.ErrInvalidCredentials:
		return errInvalidCredentials
	case auth.ErrInvalidToken:
		return errInvalidToken
	case auth.ErrTokenExpired:
		return errTokenExpired
	case auth.ErrUnauthorized:
		return errUnauthorized
	default:
		return err
	}
}

var NotFound = []error{
	errInvalidCredentials,
	errInvalidToken,
}
