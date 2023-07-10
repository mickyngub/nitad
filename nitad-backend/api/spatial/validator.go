package spatial

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditSpatialValidator(c *fiber.Ctx) error {
	sr := new(SpatialRequest)

	if err := c.BodyParser(sr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*sr)
	if err != nil {
		return errors.Throw(c, err)
	}

	spatial := new(Spatial)
	err = utils.CopyStruct(sr, spatial)
	if err != nil {
		return errors.Throw(c, err)
	}

	c.Locals("spatialBody", spatial)

	return c.Next()
}
