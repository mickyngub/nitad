package subcategory

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	subcategoryRoute fiber.Router,
) *Controller {

	controller := &Controller{service}

	subcategoryRoute.Get("/", controller.ListSubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.GetSubcategory)

	// subcategoryRoute.Use(admin.IsAuth())
	subcategoryRoute.Post("/", AddAndEditSubcategoryValidator, controller.AddSubcategory)
	subcategoryRoute.Put("/:subcategoryId", AddAndEditSubcategoryValidator, controller.EditSubcategory)
	subcategoryRoute.Delete("/:subcategoryId", controller.DeleteSubcategory)

	return controller
}

type Controller struct {
	service Service
}

// list all subcategories
func (c *Controller) ListSubcategory(ctx *fiber.Ctx) error {
	subcategories, err := c.service.ListSubcategory(ctx.Context())
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcategories})
}

// get subcategory by id
func (c *Controller) GetSubcategory(ctx *fiber.Ctx) error {
	subcategoryId := ctx.Params("subcategoryId")

	objectId, err := utils.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	subcate, err := c.service.GetSubcategoryById(ctx.Context(), objectId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcate})
}

// add a subcategory
func (c *Controller) AddSubcategory(ctx *fiber.Ctx) error {
	subcategory, ok := ctx.Locals("subcategoryBody").(*Subcategory)
	if !ok {
		return errors.Throw(ctx, errors.NewInternalServerError("Edit subcategory went wrong!"))
	}

	files, err := utils.ExtractFiles(ctx, "image")
	if err != nil {
		return errors.Throw(ctx, err)
	}

	addedSubcate, err := c.service.AddSubcategory(ctx.Context(), files, subcategory)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": addedSubcate})
}

// // edit the subcategory
func (c *Controller) EditSubcategory(ctx *fiber.Ctx) error {
	subcategoryId := ctx.Params("subcategoryId")
	objectId, err := utils.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	subcate, ok := ctx.Locals("subcategoryBody").(*Subcategory)
	if !ok {
		return errors.Throw(ctx, errors.NewInternalServerError("Edit subcategory went wrong!"))
	}
	subcate.ID = objectId

	editedSubcate, err := c.service.EditSubcategory(ctx, subcate)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": editedSubcate})
}

// delete the subcategory
func (c *Controller) DeleteSubcategory(ctx *fiber.Ctx) error {
	subcategoryId := ctx.Params("subcategoryId")
	objectId, err := utils.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	err = c.service.DeleteSubcategory(ctx.Context(), objectId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete subcategory " + subcategoryId + " successfully!"})
}
