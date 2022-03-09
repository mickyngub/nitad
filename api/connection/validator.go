package connection

import (
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditSubcategoryValidator(ctx *fiber.Ctx) error {
	subcategoryDTO := new(subcategory.SubcategoryDTO)
	if err := ctx.BodyParser(subcategoryDTO); err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*subcategoryDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Next()
}
