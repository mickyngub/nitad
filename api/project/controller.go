package project

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	projectRoute fiber.Router,
) {

	controller := &Controller{}

	projectRoute.Post("/", controller.AddProject)

}

type Controller struct {
	// service Service
}

var collectionName = database.COLLECTIONS["PROJECT"]

func (contc *Controller) AddProject(c *fiber.Ctx) error {
	p := new(ProjectRequest)
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
