package backend

import "fmt"

type HttpStatusCode int16

const (
	InternalError        HttpStatusCode = 500
	NotFoundError        HttpStatusCode = 404
	ForbiddenError       HttpStatusCode = 403
	UnauthenticatedError HttpStatusCode = 401
	BadRequestError      HttpStatusCode = 400
)

type Error struct {
	code    HttpStatusCode
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %v", e.code, e.message)
}

func NewError(code HttpStatusCode, message string) Error {
	return Error{
		code:    code,
		message: message,
	}
}
