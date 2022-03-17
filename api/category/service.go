package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListCategory(ctx context.Context) ([]*Category, errors.CustomError)
	GetCategoryById(ctx context.Context, id string) (*Category, errors.CustomError)
	AddCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError)
	EditCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError)
	DeleteCategory(ctx context.Context, id string) errors.CustomError

	SearchCategory(ctx context.Context) ([]CategorySearch, errors.CustomError)
	FindByIds2(ctx context.Context, cids []string) ([]Category, []primitive.ObjectID, errors.CustomError)

	FilterCatesWithSubcates(categories []Category, subcategories []subcategory.Subcategory) ([]Category, errors.CustomError)
	FilterCatesWithSids(categories []Category, sids []primitive.ObjectID) ([]Category, errors.CustomError)
	BindSubcategory(ctx context.Context, coid primitive.ObjectID, soid primitive.ObjectID) errors.CustomError
	UnbindSubcategory(ctx context.Context, coid primitive.ObjectID, soid primitive.ObjectID) errors.CustomError
	AddSubcategory(ctx context.Context, cid string, sid string) (*CategoryDTO, errors.CustomError)

	UpdateProjectCount(ctx context.Context, cate *Category, val int) errors.CustomError
}

type categoryService struct {
	repository         Repository
	subcategoryService subcategory.Service
}

func NewService(repository Repository, subcategoryService subcategory.Service) Service {
	return &categoryService{repository, subcategoryService}
}

func (c *categoryService) ListCategory(ctx context.Context) ([]*Category, errors.CustomError) {
	cates, err := c.repository.ListCategory(ctx)
	if err != nil {
		return nil, err
	}
	for _, cate := range cates {
		for _, subcate := range cate.Subcategory {
			subcate.Image = gcp.GetURL(subcate.Image)
		}
	}
	return cates, nil
}

// this function is reused
func (c *categoryService) GetCategoryById(ctx context.Context, id string) (*Category, errors.CustomError) {
	oid, err := database.ExtractOID(id)
	if err != nil {
		return nil, err
	}

	cate, err := c.repository.GetCategoryById(ctx, oid)
	if err != nil {
		return nil, err
	}

	return cate, nil
}

func (c *categoryService) AddCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError) {
	cateDTO.Subcategory = utils.RemoveDuplicateIds(cateDTO.Subcategory)
	subcategories, _, err := c.subcategoryService.FindByIds3(ctx, cateDTO.Subcategory)
	if err != nil {
		return cateDTO, err
	}

	result := new(CategoryDTO)
	err = database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		var txErr errors.CustomError
		result, txErr = c.addCategoryAndBindSubcategory(sessionContext, cateDTO, subcategories)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (c *categoryService) addCategoryAndBindSubcategory(ctx context.Context, cateDTO *CategoryDTO, subcategories []subcategory.Subcategory) (*CategoryDTO, errors.CustomError) {
	cateDTO, err := c.repository.AddCategory(ctx, cateDTO)
	if err != nil {
		return cateDTO, err
	}

	for _, subcate := range subcategories {
		_, err = c.subcategoryService.InsertToCategory(ctx, &subcate, cateDTO.ID)
		if err != nil {
			return cateDTO, err
		}
	}
	return cateDTO, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (c *categoryService) EditCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError) {
	oldCate, err := c.GetCategoryById(ctx, cateDTO.ID.Hex())
	if err != nil {
		return cateDTO, err
	}

	cateDTO.Subcategory = utils.RemoveDuplicateIds(cateDTO.Subcategory)
	// check if parse subcategoryIds exist
	subcategories, _, err := c.subcategoryService.FindByIds3(ctx, cateDTO.Subcategory)
	if err != nil {
		return cateDTO, err
	}

	// find the remove subcategories by checking
	// the non-intersecting of oldSubcategories and updatedSubcategories
	removeSubcategories := []*subcategory.Subcategory{}
	for _, oldSubcategory := range oldCate.Subcategory {
		if !contains(cateDTO.Subcategory, oldSubcategory.ID.Hex()) {
			removeSubcategories = append(removeSubcategories, oldSubcategory)
		}
	}

	result := new(CategoryDTO)
	err = database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		var txErr errors.CustomError
		result, txErr = c.editCategoryAndBindSubcategory(sessionContext, cateDTO, subcategories, removeSubcategories)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *categoryService) editCategoryAndBindSubcategory(ctx context.Context, cateDTO *CategoryDTO, subcategories []subcategory.Subcategory, removeSubcategories []*subcategory.Subcategory) (*CategoryDTO, errors.CustomError) {
	cateDTO, err := c.repository.EditCategory(ctx, cateDTO)
	if err != nil {
		return cateDTO, err
	}

	// set categoryId to the updated ones
	for _, subcate := range subcategories {
		_, err = c.subcategoryService.InsertToCategory(ctx, &subcate, cateDTO.ID)
		if err != nil {
			return cateDTO, err
		}
	}

	// unset the remove ones
	for _, removeSubcate := range removeSubcategories {
		_, err = c.subcategoryService.InsertToCategory(ctx, removeSubcate, primitive.NilObjectID)
		if err != nil {
			return cateDTO, err
		}
	}

	return cateDTO, nil

}

func (c *categoryService) DeleteCategory(ctx context.Context, id string) errors.CustomError {
	cate, err := c.GetCategoryById(ctx, id)
	if err != nil {
		return err
	}

	if cate.ProjectCount > 0 {
		return errors.NewBadRequestError("Cannot delete: Category " + cate.Title + " is still being used in some projects.")
	}

	return database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		return c.deleteCategory(sessionContext, cate)
	})
}

func (c *categoryService) deleteCategory(ctx context.Context, cate *Category) errors.CustomError {
	// unbind subcategory
	for _, subcate := range cate.Subcategory {
		c.subcategoryService.InsertToCategory(ctx, subcate, primitive.NilObjectID)
	}

	return c.repository.DeleteCategory(ctx, cate.ID)

}

func (c *categoryService) SearchCategory(ctx context.Context) ([]CategorySearch, errors.CustomError) {
	return c.repository.SearchCategory(ctx)
}

func (c *categoryService) AddSubcategory(ctx context.Context, cid string, sid string) (*CategoryDTO, errors.CustomError) {
	coid, err := database.ExtractOID(cid)
	if err != nil {
		return nil, err
	}

	cateDTO, err := c.repository.GetCategoryByIdNoLookup(ctx, coid)
	if err != nil {
		return nil, err
	}

	// find new added sid
	addedSubcate, err := c.subcategoryService.GetSubcategoryById(ctx, sid)
	if err != nil {
		return nil, err
	}

	// append new one to the existing one
	// cate := new(Category)
	// utils.CopyStruct(cateDTO, cate)
	cateDTO.Subcategory = append(cateDTO.Subcategory, sid)
	cateDTO.Subcategory = utils.RemoveDuplicateIds(cateDTO.Subcategory)
	cateDTO, err = c.repository.EditCategory(ctx, cateDTO)

	if err != nil {
		return cateDTO, err
	}

	// update subcate bounded to cate
	_, err = c.subcategoryService.InsertToCategory(ctx, addedSubcate, cateDTO.ID)
	if err != nil {
		return cateDTO, err
	}

	return cateDTO, err

}

func (c *categoryService) UpdateProjectCount(ctx context.Context, cate *Category, val int) errors.CustomError {
	return c.repository.UpdateProjectCount(ctx, cate.ID, val)
}
