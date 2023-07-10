package errors

import "github.com/gofiber/fiber/v2"

type cacheError struct {
	message string
}

func NewCacheError(message string) CustomError {
	return &cacheError{message}
}

func (b *cacheError) Code() int {
	return fiber.StatusInternalServerError
}

func (b *cacheError) Error() string {
	return b.message
}
