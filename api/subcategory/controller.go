package subcategory

import (
	"log"

	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	subcategoryRoute fiber.Router,
) {

	controller := &Controller{}

	subcategoryRoute.Get("/", controller.ListSubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.GetSubcategory)

	subcategoryRoute.Use(admin.IsAuth())
	subcategoryRoute.Post("/", AddAndEditSubcategoryValidator, controller.AddSubcategory)
	subcategoryRoute.Put("/:subcategoryId", AddAndEditSubcategoryValidator, controller.EditSubcategory)
	subcategoryRoute.Delete("/:subcategoryId", controller.DeleteSubcategory)
}

type Controller struct{}

// list all subcategories
func (contc *Controller) ListSubcategory(c *fiber.Ctx) error {
	subcategories, err := FindAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcategories})
}

// get subcategory by id
func (contc *Controller) GetSubcategory(c *fiber.Ctx) error {
	subcategoryId := c.Params("subcategoryId")

	objectId, err := utils.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var result Subcategory
	if result, err = GetById(objectId); err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a subcategory
func (contc *Controller) AddSubcategory(c *fiber.Ctx) error {
	files, err := utils.ExtractFiles(c, "image")
	if err != nil {
		return errors.Throw(c, err)
	}

	imageFilename, err := gcp.UploadFile(c.Context(), files[0], collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}

	sr := new(SubcategoryRequest)
	c.BodyParser(sr)
	var subcategory Subcategory
	subcategory.Title = sr.Title
	subcategory.Image = imageFilename

	result, err := Add(&subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// // edit the subcategory
func (contc *Controller) EditSubcategory(c *fiber.Ctx) error {
	subcategoryId := c.Params("subcategoryId")
	objectId, err := utils.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	updateSubcategory, ok := c.Locals("subcategoryBody").(*Subcategory)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Edit subcategory went wrong!"))
	}

	updateSubcategory.ID = objectId
	updateSubcategory, err = HandleUpdateImage(c, updateSubcategory)
	log.Println(updateSubcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	result, err := Edit(updateSubcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// delete the subcategory
func (cont *Controller) DeleteSubcategory(c *fiber.Ctx) error {
	subcategoryId := c.Params("subcategoryId")
	objectId, err := utils.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	oldSubcategory, err := GetById(objectId)
	if err != nil {
		return err
	}

	gcp.DeleteFile(c.Context(), oldSubcategory.Image, collectionName)

	err = Delete(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete subcategory " + subcategoryId + " successfully!"})
}
