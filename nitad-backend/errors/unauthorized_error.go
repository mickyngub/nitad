package errors

import "github.com/gofiber/fiber/v2"

type unAuthorizedError struct {
	message string
}

var (
	NewExpiredToken       = NewUnAuthorizedError("token is expired")
	NewInvalidToken       = NewUnAuthorizedError("token is invalid")
	NewInvalidCredentials = NewUnAuthorizedError("invalid credentials")
)

func NewUnAuthorizedError(message string) CustomError {
	return &unAuthorizedError{message}
}

func (b *unAuthorizedError) Code() int {
	return fiber.StatusUnauthorized
}

func (b *unAuthorizedError) Error() string {
	return b.message
}
