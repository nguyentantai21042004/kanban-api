package labels

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrLabelNotFound = errors.New("label not found")
)
