package category

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditCategoryValidator(ctx *fiber.Ctx) error {
	cateDTO := new(CategoryDTO)

	if err := ctx.BodyParser(cateDTO); err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*cateDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Next()
}
