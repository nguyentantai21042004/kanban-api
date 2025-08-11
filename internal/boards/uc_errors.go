package boards

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrNotFound      = errors.New("board not found")
)
