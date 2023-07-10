package connection

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	AddSubcategory(ctx context.Context, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError)
	EditSubcategory(ctx context.Context, id string, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError)
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

func (c *connectionService) AddSubcategory(ctx context.Context, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError) {

	addedSubcate := new(subcategory.Subcategory)
	err := database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		var txErr errors.CustomError
		addedSubcate, txErr = c.addSubcategory(sessionContext, subcateDTO)
		return txErr
	})
	if err != nil {
		return nil, err
	}

	return addedSubcate, nil

}

func (c *connectionService) addSubcategory(ctx context.Context, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError) {
	addedSubcate, err := c.subcategoryService.AddSubcategory(ctx, subcateDTO)
	if err != nil {
		return nil, err
	}

	if subcateDTO.CategoryId != primitive.NilObjectID {
		// if there is parse categoryId, check if the cate exists
		cate, err := c.categoryService.GetCategoryById(ctx, subcateDTO.CategoryId.Hex())
		if err != nil {
			return nil, err
		}
		err = c.categoryService.BindSubcategory(ctx, cate.ID, addedSubcate.ID)
		if err != nil {
			return nil, err
		}
	}
	return addedSubcate, nil

}

func (c *connectionService) EditSubcategory(ctx context.Context, id string, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError) {
	editedSubcate := new(subcategory.Subcategory)
	err := database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		var txErr errors.CustomError
		editedSubcate, txErr = c.editSubcategory(sessionContext, id, subcateDTO)
		return txErr
	})
	if err != nil {
		return nil, err
	}

	return editedSubcate, nil
}

func (c *connectionService) editSubcategory(ctx context.Context, id string, subcateDTO *subcategory.SubcategoryDTO) (*subcategory.Subcategory, errors.CustomError) {
	oldSubcate, err := c.subcategoryService.GetSubcategoryById(ctx, id)
	if err != nil {
		return nil, err
	}

	if oldSubcate.CategoryId != primitive.NilObjectID {
		// meaning that there was category binded to
		// need to unbind from category
		err = c.categoryService.UnbindSubcategory(ctx, oldSubcate.CategoryId, oldSubcate.ID)
		if err != nil {
			return nil, err
		}
	}

	subcateDTO.ID = oldSubcate.ID
	editedSubcate, err := c.subcategoryService.EditSubcategory(ctx, subcateDTO)
	if err != nil {
		return nil, err
	}

	err = c.categoryService.BindSubcategory(ctx, subcateDTO.CategoryId, editedSubcate.ID)
	if err != nil {
		return nil, err
	}
	return editedSubcate, nil
}
