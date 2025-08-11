package repository

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrFieldRequired = errors.New("field required")
)
