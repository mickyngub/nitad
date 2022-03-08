package connection

import (
	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	connectionRoute fiber.Router,
) *Controller {

	controller := &Controller{service}

	connectionRoute.Put("/subcategory", controller.EditSubcategory)

	connectionRoute.Use(admin.IsAuth())
	connectionRoute.Post("/subcategory", AddAndEditSubcategoryValidator, controller.AddSubcategory)
	connectionRoute.Put("/subcategory/:subcategoryId", AddAndEditSubcategoryValidator, controller.EditSubcategory)

	return controller
}

type Controller struct {
	service Service
}

func (c *Controller) AddSubcategory(ctx *fiber.Ctx) error {
	subcategoryDTO := new(subcategory.SubcategoryDTO)
	ctx.BodyParser(subcategoryDTO)

	files, err := utils.ExtractFiles(ctx, "image")
	if err != nil {
		return errors.Throw(ctx, err)
	}
	subcategoryDTO.Image = files[0]

	addedSubcate, err := c.service.AddSubcategory(ctx, subcategoryDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": addedSubcate})

}

func (c *Controller) EditSubcategory(ctx *fiber.Ctx) error {
	subcategoryId := ctx.Params("subcategoryId")

	subcategoryDTO := new(subcategory.SubcategoryDTO)
	ctx.BodyParser(subcategoryDTO)

	images, err := utils.ExtractUpdatedFiles(ctx, "images")
	if err != nil {
		return errors.Throw(ctx, err)
	}

	if images == nil {
		subcategoryDTO.Image = nil
	} else {
		subcategoryDTO.Image = images[0]
	}

	editedSubcate, err := c.service.EditSubcategory(ctx, subcategoryId, subcategoryDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": editedSubcate})

}
