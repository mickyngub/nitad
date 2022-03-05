package connection

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	connectionRoute fiber.Router,
) *Controller {

	controller := &Controller{}

	connectionRoute.Get("/unsetSubcategory", controller.ListUnsetSubcategory)

	return controller
}

type Controller struct{}

func (contc *Controller) ListUnsetSubcategory(c *fiber.Ctx) error {
	subcategories, err := GetSubcategoryThatAreNotInAnyCategory()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcategories})
}
