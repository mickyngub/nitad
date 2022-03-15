package errors

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CustomError interface {
	Code() int
	Error() string
}

func Throw(c *fiber.Ctx, ce CustomError) error {
	zap.S().Warn("Error Throw: ", ce.Code(), ": ", ce.Error())
	return c.Status(ce.Code()).JSON(fiber.Map{"success": false, "result": ce.Error()})
}
