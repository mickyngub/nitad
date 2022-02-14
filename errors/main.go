package errors

import "github.com/gofiber/fiber/v2"

type CustomError interface {
	Code() int
	Error() string
}

func Throw(c *fiber.Ctx, ce CustomError) error {
	return c.Status(ce.Code()).JSON(fiber.Map{"success": false, "result": ce.Error()})
}
