package errors

import "github.com/gofiber/fiber/v2"

type internalServerError struct {
	message string
}

// NewinternalServerError used when the error caused from
// any wrong logic or unexpected error in the application
func NewInternalServerError(message string) CustomError {
	return &internalServerError{message}
}

func (b *internalServerError) Code() int {
	return fiber.StatusInternalServerError
}

func (b *internalServerError) Error() string {
	return "500: " + b.message
}
