package response

import (
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int
	Err  string
}

func (e *Error) Error() string {
	return e.Err
}

func New(code int, err string) error {
	return &Error{Code: code, Err: err}
}

var (
	ErrBadRequest          = New(fiber.StatusBadRequest, "Bad Request")
	ErrForeignKeyViolation = New(fiber.StatusForbidden, "Foreign Key Violation")
)
