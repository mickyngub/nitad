package subcategory

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewController(
	subcategoryRoute fiber.Router,
) {

	controller := &Controller{}

	subcategoryRoute.Get("/", controller.Listsubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.Getsubcategory)

	//TODO add AUTH for POST/PUT/DELETE

	subcategoryRoute.Post("/", controller.Addsubcategory)
	// subcategoryRoute.Put("/:subcategoryId", controller.editsubcategory)
	// subcategoryRoute.Delete("/:subcategoryId", controller.deletesubcategory)
}

type Controller struct {
	// service Service
}

var collectionName = database.COLLECTIONS["SUBCATEGORY"]

// get subcategory by id
func (contc *Controller) Getsubcategory(c *fiber.Ctx) error {
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

// list all subcategories
func (contc *Controller) Listsubcategory(c *fiber.Ctx) error {
	subcategories, err := FindAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcategories})
}

// add a subcategory
func (contc *Controller) Addsubcategory(c *fiber.Ctx) error {

	p := new(Subcategory)

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

// // edit the subcategory
// func (contc *Controller) editsubcategory(c *fiber.Ctx) error {}

// // delete the subcategory
// func (cont *Controller) deletesubcategory(c *fiber.Ctx) error {}
