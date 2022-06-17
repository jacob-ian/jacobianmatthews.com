package backend

import "errors"

var (
	InternalErrorMsg = "internal_error"
	InternalError    = errors.New(InternalErrorMsg)
)
