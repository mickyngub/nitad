package subcategory

import (
	"fmt"

	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditSubcategoryValidator(c *fiber.Ctx) error {
	subcategoryDTO := new(SubcategoryDTO)
	if err := c.BodyParser(subcategoryDTO); err != nil {
		fmt.Println("error", err.Error())
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*subcategoryDTO)
	if err != nil {
		return errors.Throw(c, err)
	}

	subcategory := new(Subcategory)
	err = utils.CopyStruct(subcategoryDTO, subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Next()
}
