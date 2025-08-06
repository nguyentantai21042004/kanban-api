package lists

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrNotFound      = errors.New("not found")
)
