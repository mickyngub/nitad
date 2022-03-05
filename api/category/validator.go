package category

import (
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditCategoryValidator(ctx *fiber.Ctx) error {
	cr := new(CategoryRequest)

	if err := ctx.BodyParser(cr); err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*cr)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	category := new(Category)
	subcategories, sids, err := subcategory.FindByIds(cr.Subcategory)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	category.Subcategory = subcategories
	category.Title = cr.Title
	ctx.Locals("categoryBody", category)
	ctx.Locals("sids", sids)

	return ctx.Next()
}
