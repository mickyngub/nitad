package subcategory

import (
	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	subcategoryRoute fiber.Router,
) *Controller {

	controller := &Controller{service}

	subcategoryRoute.Get("/", controller.ListSubcategory)
	subcategoryRoute.Get("/unset", controller.ListUnsetSubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.GetSubcategoryById)

	subcategoryRoute.Use(admin.IsAuth())
	// subcategoryRoute.Post("/", AddAndEditSubcategoryValidator, controller.AddSubcategory)
	// subcategoryRoute.Put("/:subcategoryId", AddAndEditSubcategoryValidator, controller.EditSubcategory)
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

// list all unset subcategories
func (c *Controller) ListUnsetSubcategory(ctx *fiber.Ctx) error {
	unsetSubcategories, err := c.service.ListUnsetSubcategory(ctx.Context())
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": unsetSubcategories})
}

// get subcategory by id
func (c *Controller) GetSubcategoryById(ctx *fiber.Ctx) error {
	subcategoryId := ctx.Params("subcategoryId")

	subcate, err := c.service.GetSubcategoryById(ctx.Context(), subcategoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcate})
}

// add a subcategory
// func (c *Controller) AddSubcategory(ctx *fiber.Ctx) error {
// 	subcategoryDTO := new(SubcategoryDTO)
// 	ctx.BodyParser(subcategoryDTO)

// 	files, err := utils.ExtractFiles(ctx, "image")
// 	if err != nil {
// 		return errors.Throw(ctx, err)
// 	}
// 	subcategoryDTO.Image = files[0]

// 	addedSubcate, err := c.service.AddSubcategory(ctx.Context(), subcategoryDTO)
// 	if err != nil {
// 		return errors.Throw(ctx, err)
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": addedSubcate})
// }

// // edit the subcategory
// func (c *Controller) EditSubcategory(ctx *fiber.Ctx) error {
// 	subcategoryId := ctx.Params("subcategoryId")
// 	objectId, err := utils.IsValidObjectId(subcategoryId)
// 	if err != nil {
// 		return errors.Throw(ctx, err)
// 	}

// 	subcategoryDTO := new(SubcategoryDTO)
// 	ctx.BodyParser(subcategoryDTO)
// 	subcategoryDTO.ID = objectId

// 	images, err := utils.ExtractUpdatedFiles(ctx, "images")
// 	if err != nil {
// 		return errors.Throw(ctx, err)
// 	}

// 	if images == nil {
// 		subcategoryDTO.Image = nil
// 	} else {
// 		subcategoryDTO.Image = images[0]
// 	}

// 	editedSubcate, err := c.service.EditSubcategory(ctx, subcategoryDTO)
// 	if err != nil {
// 		return errors.Throw(ctx, err)
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": editedSubcate})
// }

// delete the subcategory
func (c *Controller) DeleteSubcategory(ctx *fiber.Ctx) error {
	subcategoryId := ctx.Params("subcategoryId")

	err := c.service.DeleteSubcategory(ctx.Context(), subcategoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete subcategory " + subcategoryId + " successfully!"})
}
