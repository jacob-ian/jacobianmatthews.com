package backend

import "fmt"

const (
	InternalError        int = 500
	NotImplemented       int = 501
	NotFoundError        int = 404
	ForbiddenError       int = 403
	UnauthenticatedError int = 401
	BadRequestError      int = 400
)

type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.code, e.message)
}

func (e Error) GetCode() int {
	return e.code
}

func (e Error) GetMessage() string {
	return e.message
}

func (e Error) Is(code int) bool {
	return e.code == code
}

func NewError(code int, message string) Error {
	return Error{
		code:    code,
		message: message,
	}
}
