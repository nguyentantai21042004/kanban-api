package lists

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrListNotFound  = errors.New("list not found")
)
