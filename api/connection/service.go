package connection

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	AddSubcategory(ctx *fiber.Ctx, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError)
	EditSubcategory(ctx *fiber.Ctx, id string, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError)
}

type connectionService struct {
	subcategoryService subcategory.Service
	categoryService    category.Service
}

func NewService(subcategoryService subcategory.Service, categoryService category.Service) Service {
	return &connectionService{
		subcategoryService,
		categoryService,
	}
}

func (c *connectionService) AddSubcategory(ctx *fiber.Ctx, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError) {

	addedSubcate, err := c.subcategoryService.AddSubcategory(ctx.Context(), subcateDTO)
	if err != nil {
		return nil, err
	}

	if subcateDTO.CategoryId != primitive.NilObjectID {
		// if there is parse categoryId, check if the cate exists
		cate, err := c.categoryService.GetCategoryById(ctx.Context(), subcateDTO.CategoryId.Hex())
		if err != nil {
			return nil, err
		}
		err = c.categoryService.BindSubcategory(ctx.Context(), cate.ID, addedSubcate.ID)
		if err != nil {
			return nil, err
		}
	}

	return addedSubcate, nil

}

func (c *connectionService) EditSubcategory(ctx *fiber.Ctx, id string, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError) {
	oldSubcate, err := c.subcategoryService.GetSubcategoryById(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	//TODO: tx these two service in case one is broken -> revert all
	if oldSubcate.CategoryId != primitive.NilObjectID {
		// meaning that there was category binded to
		// need to unbind from category
		err = c.categoryService.UnbindSubcategory(ctx.Context(), oldSubcate.CategoryId, oldSubcate.ID)
		if err != nil {
			return nil, err
		}
	}

	subcateDTO.ID = oldSubcate.ID
	editedSubcate, err := c.subcategoryService.EditSubcategory(ctx, subcateDTO)
	if err != nil {
		return nil, err
	}

	err = c.categoryService.BindSubcategory(ctx.Context(), subcateDTO.CategoryId, editedSubcate.ID)
	if err != nil {
		return nil, err
	}

	return editedSubcate, nil
}
