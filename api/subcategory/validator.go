package subcategory

import (
	"fmt"

	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditSubcategoryValidator(ctx *fiber.Ctx) error {
	subcategoryDTO := new(SubcategoryDTO)
	if err := ctx.BodyParser(subcategoryDTO); err != nil {
		fmt.Println("error", err.Error())
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*subcategoryDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	subcategory := new(Subcategory)
	err = utils.CopyStruct(subcategoryDTO, subcategory)
	if err != nil {
		return errors.Throw(ctx, err)
	}
	ctx.Locals("subcategoryBody", subcategory)

	return ctx.Next()
}
