package errors

import "github.com/gofiber/fiber/v2"

type badRequestError struct {
	message string
}

// NewBadRequestError used when client send wrong type of input
// e.g. expect application/json but got plain/text
// 			or wrong username/password
func NewBadRequestError(message string) CustomError {
	return &badRequestError{message}
}

func (b *badRequestError) Code() int {
	return fiber.StatusBadRequest
}

func (b *badRequestError) Error() string {
	return b.message
}
