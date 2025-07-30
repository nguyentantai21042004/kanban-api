package role

import "errors"

var (
	ErrRoleNotFound  = errors.New("role not found")
	ErrFieldRequired = errors.New("field required")
)
