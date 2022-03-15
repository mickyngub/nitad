package errors

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CustomError interface {
	Code() int
	Error() string
}

func Throw(ctx *fiber.Ctx, ce CustomError) error {
	zap.S().Warn("Error Throw: ", ctx.Path(), ": ", ce.Code(), ": ", ce.Error())
	return ctx.Status(ce.Code()).JSON(fiber.Map{"success": false, "result": ce.Error()})
}
