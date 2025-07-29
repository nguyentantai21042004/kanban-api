package cards

import "errors"

var (
	ErrFieldRequired = errors.New("field required")
	ErrCardNotFound  = errors.New("card not found")
	ErrListNotFound  = errors.New("list not found")
)
