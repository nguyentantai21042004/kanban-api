package http

import (
	"github.com/nguyentantai21042004/kanban-api/internal/role"
	pkgErrors "github.com/nguyentantai21042004/kanban-api/pkg/errors"
)

var (
	errWrongQuery    = pkgErrors.NewHTTPError(10401, "Wrong query")
	errNotFound      = pkgErrors.NewHTTPError(10403, "Role not found")
	errFieldRequired = pkgErrors.NewHTTPError(10404, "Field required")
)

func (h handler) mapErrorCode(err error) error {
	switch err {
	case role.ErrRoleNotFound:
		return errNotFound
	case role.ErrFieldRequired:
		return errFieldRequired
	default:
		return err
	}
}

var NotFound = []error{
	errNotFound,
}
