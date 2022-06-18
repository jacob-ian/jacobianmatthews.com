package backend

import "errors"

var (
	InternalErrorMsg = "internal_error"
	InternalError    = errors.New(InternalErrorMsg)
	NotFoundErrorMsg = "not_found"
	NotFoundError    = errors.New(NotFoundErrorMsg)
)
