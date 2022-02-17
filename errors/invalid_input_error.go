package errors

import "github.com/gofiber/fiber/v2"

type invalidInputError struct {
	fields []string
}

func NewInvalidInputError(fields []string) CustomError {
	return &invalidInputError{fields}
}

func (b *invalidInputError) Code() int {
	return fiber.StatusBadRequest
}

func (b *invalidInputError) Error() string {
	str := ""
	for _, field := range b.fields {
		str += field + ", "
	}
	return "Invalid input: " + str
}
