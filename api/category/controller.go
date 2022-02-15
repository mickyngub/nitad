package category

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewController(
	categoryRoute fiber.Router,
) {

	controller := &Controller{}

	categoryRoute.Get("/", controller.ListCategory)
	categoryRoute.Get("/:categoryId", controller.GetCategory)

	//TODO add AUTH for POST/PUT/DELETE

	categoryRoute.Post("/", controller.AddCategory)
	// categoryRoute.Put("/:categoryId", controller.editcategory)
	// categoryRoute.Delete("/:categoryId", controller.deletecategory)
}

type Controller struct {
	// service Service
}

var collectionName = database.COLLECTIONS["CATEGORY"]

// list all categories
func (contc *Controller) ListCategory(c *fiber.Ctx) error {
	categories, err := FindAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": categories})
}

// get category by id
func (contc *Controller) GetCategory(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")

	objectId, err := functions.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var result bson.M
	if result, err = FindById(objectId); err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a category
func (contc *Controller) AddCategory(c *fiber.Ctx) error {

	p := new(CategoryRequest)
	if err := c.BodyParser(p); err != nil {
		return errors.Throw(c, errors.InvalidInput)
	}

	result, err := Add(p)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

// // edit the category
// func (contc *Controller) editCategory(c *fiber.Ctx) error {}

// // delete the category
// func (cont *Controller) deleteCategory(c *fiber.Ctx) error {}
