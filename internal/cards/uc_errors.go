package cards

import "errors"

var (
	ErrFieldRequired         = errors.New("field required")
	ErrCardNotFound          = errors.New("card not found")
	ErrListNotFound          = errors.New("list not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrAttachmentNotFound    = errors.New("attachment not found")
	ErrTagNotFound           = errors.New("tag not found")
	ErrInvalidTimeRange      = errors.New("invalid time range")
	ErrChecklistItemNotFound = errors.New("checklist item not found")
)
