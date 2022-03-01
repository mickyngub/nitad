package search

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	searchRoute fiber.Router,
) {

	controller := &Controller{}

	searchRoute.Get("/", controller.SearchAll)
}

type Controller struct{}

// list all neccessary components: project/ category/ subcategory
func (contc *Controller) SearchAll(c *fiber.Ctx) error {
	result, err := SearchAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}
