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

	categoryRoute.Get("/", controller.Listcategory)
	categoryRoute.Get("/:categoryId", controller.Getcategory)

	//TODO add AUTH for POST/PUT/DELETE

	categoryRoute.Post("/", controller.Addcategory)
	// categoryRoute.Put("/:categoryId", controller.editcategory)
	// categoryRoute.Delete("/:categoryId", controller.deletecategory)
}

type Controller struct {
	// service Service
}

var collectionName = database.COLLECTIONS["CATEGORY"]

// get category by id
func (contc *Controller) Getcategory(c *fiber.Ctx) error {
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

// list all categories
func (contc *Controller) Listcategory(c *fiber.Ctx) error {
	categories, err := FindAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": categories})
}

// add a category
func (contc *Controller) Addcategory(c *fiber.Ctx) error {

	p := new(CategoryRequest)
	//TODO: handle this bodyParser middleware
	if err := c.BodyParser(p); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	result, err := Add(p)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

// // edit the category
// func (contc *Controller) editcategory(c *fiber.Ctx) error {}

// // delete the category
// func (cont *Controller) deletecategory(c *fiber.Ctx) error {}
