package search

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	searchRoute fiber.Router,
) {
	controller := &Controller{}
	searchRoute.Get("/", controller.SearchAll)
}

type Controller struct {
	service Service
}

// list all neccessary components: project/ category/ subcategory
func (c *Controller) SearchAll(ctx *fiber.Ctx) error {
	searchResult, err := c.service.SearchAll(ctx)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": searchResult})
}
