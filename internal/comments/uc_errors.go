package comments

import "errors"

var (
	ErrFieldRequired         = errors.New("field required")
	ErrCommentNotFound       = errors.New("comment not found")
	ErrCardNotFound          = errors.New("card not found")
	ErrParentCommentNotFound = errors.New("parent comment not found")
)
