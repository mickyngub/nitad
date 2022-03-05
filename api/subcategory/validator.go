package subcategory

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditSubcategoryValidator(c *fiber.Ctx) error {
	sr := new(SubcategoryRequest)

	if err := c.BodyParser(sr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError("EIEI "+err.Error()))
	}

	err := validators.ValidateStruct(*sr)
	if err != nil {
		return errors.Throw(c, err)
	}

	subcategory := new(Subcategory)
	err = utils.CopyStruct(sr, subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	c.Locals("subcategoryBody", subcategory)

	return c.Next()
}
