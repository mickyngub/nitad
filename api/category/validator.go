package category

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func AddAndEditCategoryValidator(ctx *fiber.Ctx) error {
	categoryDTO := new(CategoryDTO)

	if err := ctx.BodyParser(categoryDTO); err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*categoryDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	// osids := []primitive.ObjectID{}
	// for _, sid := range categoryDTO.Subcategory {
	// 	osid, err := utils.IsValidObjectId(sid)
	// 	if err != nil {
	// 		return errors.Throw(ctx, err)
	// 	}
	// 	osids = append(osids, osid)
	// }

	// cateDTO := new(CategoryDTO)
	// utils.CopyStruct(categoryDTO, cateDTO)
	// cateDTO.Subcategory = osids
	// ctx.Locals("cateDTO", cateDTO)

	return ctx.Next()
}
