package category

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func AddAndEditCategoryValidator(ctx *fiber.Ctx) error {
	cateRequest := new(CategoryRequest)

	if err := ctx.BodyParser(cateRequest); err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*cateRequest)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	osids := []primitive.ObjectID{}
	for _, sid := range cateRequest.Subcategory {
		osid, err := utils.IsValidObjectId(sid)
		zap.S().Info("hello", osid)
		if err != nil {
			return errors.Throw(ctx, err)
		}
		osids = append(osids, osid)
	}

	cateDTO := new(CategoryDTO)
	utils.CopyStruct(cateRequest, cateDTO)
	cateDTO.Subcategory = osids
	ctx.Locals("cateDTO", cateDTO)

	return ctx.Next()
}
