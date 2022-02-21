package admin

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func SignupValidator(c *fiber.Ctx) error {
	a := new(AdminSignup)

	if err := c.BodyParser(a); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*a)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Next()
}

func LoginValidator(c *fiber.Ctx) error {
	a := new(Admin)

	if err := c.BodyParser(a); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*a)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Next()
}
