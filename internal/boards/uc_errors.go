package boards

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrBoardNotFound = errors.New("board not found")
)
