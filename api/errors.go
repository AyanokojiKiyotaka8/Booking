package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *Error) Error() string {
	return e.Err
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if displayedError, ok := err.(*Error); ok {
		return c.Status(displayedError.Code).JSON(displayedError)
	}
	displayedError := ErrInternalServerError()
	return c.Status(displayedError.Code).JSON(displayedError)
}

func NewError(code int, err string) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func ErrUnauthorized() *Error {
	return &Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}

func ErrInternalServerError() *Error {
	return &Error{
		Code: http.StatusInternalServerError,
		Err:  "internal server error",
	}
}

func ErrBadRequest() *Error {
	return &Error{
		Code: http.StatusBadRequest,
		Err:  "bad request",
	}
}

func ErrNotFound() *Error {
	return &Error{
		Code: http.StatusNotFound,
		Err:  "resource not found",
	}
}

func ErrForbidden() *Error {
	return &Error{
		Code: http.StatusForbidden,
		Err:  "forbidden",
	}
}
