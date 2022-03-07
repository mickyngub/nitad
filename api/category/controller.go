package category

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	categoryRoute fiber.Router,
) {

	controller := &Controller{service}

	categoryRoute.Get("/", controller.ListCategory)
	categoryRoute.Get("/:categoryId", controller.GetCategory)

	// categoryRoute.Use(admin.IsAuth())
	categoryRoute.Post("/", AddAndEditCategoryValidator, controller.AddCategory)
	categoryRoute.Put("/:categoryId", AddAndEditCategoryValidator, controller.EditCategory)
	categoryRoute.Delete("/:categoryId", controller.DeleteCategory)
}

type Controller struct {
	service Service
}

// list all categories
func (c *Controller) ListCategory(ctx *fiber.Ctx) error {
	categories, err := c.service.ListCategory(ctx.Context())
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": categories})
}

// get category by id
func (c *Controller) GetCategory(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("categoryId")

	objectId, err := utils.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	cate, err := c.service.GetCategoryById(ctx.Context(), objectId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": cate})
}

// add a category
func (c *Controller) AddCategory(ctx *fiber.Ctx) error {
	cateDTO := new(CategoryDTO)
	ctx.BodyParser(cateDTO)

	addedCate, err := c.service.AddCategory(ctx.Context(), cateDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": addedCate})
}

// edit the category
func (c *Controller) EditCategory(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("categoryId")
	objectId, err := utils.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	cateDTO := new(CategoryDTO)
	ctx.BodyParser(cateDTO)
	cateDTO.ID = objectId

	editedCate, err := c.service.EditCategory(ctx.Context(), cateDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": editedCate})

}

// delete the category
func (c *Controller) DeleteCategory(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("categoryId")
	objectId, err := utils.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	err = c.service.DeleteCategory(ctx.Context(), objectId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete category " + categoryId + " successfully!"})
}
