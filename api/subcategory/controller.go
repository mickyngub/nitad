package subcategory

import (
	"log"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewController(
	subcategoryRoute fiber.Router,
) {

	controller := &Controller{}

	subcategoryRoute.Get("/", controller.ListSubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.GetSubcategory)

	//TODO add AUTH for POST/PUT/DELETE
	subcategoryRoute.Post("/", AddSubcategoryValidator, controller.AddSubcategory)
	subcategoryRoute.Put("/:subcategoryId", controller.EditSubcategory)
	subcategoryRoute.Delete("/:subcategoryId", controller.DeleteSubcategory)
}

type Controller struct{}

var collectionName = database.COLLECTIONS["SUBCATEGORY"]

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

	objectId, err := functions.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var result bson.M
	if result, err = FindById(objectId); err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a subcategory
func (contc *Controller) AddSubcategory(c *fiber.Ctx) error {
	files, err := functions.ExtractFiles(c, "image")
	if err != nil {
		return errors.Throw(c, err)
	}

	imageURLs, err := gcp.UploadImages(files, collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}

	sr := new(SubcategoryRequest)
	c.BodyParser(sr)
	var subcategory Subcategory
	subcategory.Title = sr.Title
	subcategory.Image = imageURLs[0]

	result, err := Add(&subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// // edit the subcategory
func (contc *Controller) EditSubcategory(c *fiber.Ctx) error {
	p := new(Subcategory)

	if err := c.BodyParser(p); err != nil {
		log.Println(err.Error())
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	subcategoryId := c.Params("subcategoryId")
	objectId, err := functions.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	p, err = HandleUpdateImage(c, p, objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	result, err := Edit(objectId, p)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// delete the subcategory
func (cont *Controller) DeleteSubcategory(c *fiber.Ctx) error {
	subcategoryId := c.Params("subcategoryId")
	objectId, err := functions.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	err = HandleDeleteImage(objectId)
	log.Println("Delete image gcp error:", err.Error())

	err = Delete(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete subcategory successfully!"})

}
