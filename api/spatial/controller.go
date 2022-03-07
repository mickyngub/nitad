package spatial

import (
	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	spatialRoute fiber.Router,
) {

	controller := &Controller{}

	spatialRoute.Get("/", controller.GetSpatial)

	spatialRoute.Use(admin.IsAuth())
	spatialRoute.Post("/", AddAndEditSpatialValidator, controller.AddSpatial)
	spatialRoute.Put("/", AddAndEditSpatialValidator, controller.EditSpatial)
	spatialRoute.Delete("/", controller.DeleteSpatial)
}

type Controller struct{}

// get that one spatial in the db
func (contc *Controller) GetSpatial(c *fiber.Ctx) error {
	spatial, err := GetOneSpatial()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": spatial})
}

// add that one spatial in the db
func (contc *Controller) AddSpatial(c *fiber.Ctx) error {
	_, err := GetOneSpatial()
	if err != nil && err.Error() != "no spatial link yet, please create one first" {
		return errors.Throw(c, err)
	}
	if err == nil { // meaning that there is spatial link passed
		return errors.Throw(c, errors.NewBadRequestError("Can only create 1 spatial link, please edit it instead"))
	}

	parseSpatial, ok := c.Locals("spatialBody").(*Spatial)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Add Spatial went wrong!"))
	}

	newSpatial, err := Add(parseSpatial)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": newSpatial})
}

// edit that one spatial in the db
func (contc *Controller) EditSpatial(c *fiber.Ctx) error {
	spatial, err := GetOneSpatial()
	if err != nil {
		return errors.Throw(c, err)
	}

	parseSpatial, ok := c.Locals("spatialBody").(*Spatial)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Edit Spatial went wrong!"))

	}

	spatial.Link = parseSpatial.Link

	updatedSpatial, err := Edit(&spatial)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": updatedSpatial})
}

// delete that one spatial in the db
func (contc *Controller) DeleteSpatial(c *fiber.Ctx) error {
	spatial, err := GetOneSpatial()
	if err != nil {
		return errors.Throw(c, err)
	}

	err = Delete(spatial.ID)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete spatial " + spatial.Link + " successfully!"})
}
