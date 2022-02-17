package category

import (
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditCategoryValidator(c *fiber.Ctx) error {
	cr := new(CategoryRequest)

	if err := c.BodyParser(cr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*cr)
	if err != nil {
		return errors.Throw(c, err)
	}

	category := new(Category)
	subcategories, sids, err := subcategory.FindByIds(cr.Subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	category.Subcategory = subcategories
	category.Title = cr.Title
	c.Locals("categoryBody", category)
	c.Locals("sids", sids)

	return c.Next()
}
