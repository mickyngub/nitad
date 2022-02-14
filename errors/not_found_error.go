package errors

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type notFoundError struct {
	resourceName string
}

// NewNotFoundError used when the resource client is asking does not exist
func NewNotFoundError(resourceName string) CustomError {
	return &notFoundError{resourceName}
}

func (n *notFoundError) Code() int {
	return fiber.StatusNotFound
}

func (n *notFoundError) Error() string {
	return fmt.Sprintf("%s not found", n.resourceName)
}

// func ThrowNotFoundError(c *fiber.Ctx, resourceName string) error {
// 	notFoundError := NewNotFoundError(resourceName)
// 	return c.Status(notFoundError.Code()).JSON(fiber.Map{"result": notFoundError.Error()})
// 	// return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"result": fmt.Sprintf("%s not found", resourceName)})
// }
