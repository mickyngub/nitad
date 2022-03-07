package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListCategory(ctx context.Context) ([]Category, errors.CustomError)
	GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	AddCategory(ctx *fiber.Ctx, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError)
	EditCategory(ctx *fiber.Ctx, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError)
	DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	SearchCategory(ctx *fiber.Ctx) ([]CategorySearch, errors.CustomError)
	FindByIds2(ctx context.Context, cids []string) ([]Category, []primitive.ObjectID, errors.CustomError)

	AddSubcategory(ctx *fiber.Ctx, oid primitive.ObjectID, sid primitive.ObjectID) (*CategoryDTO, errors.CustomError)
}

type categoryService struct {
	repository         Repository
	subcategoryService subcategory.Service
}

func NewService(repository Repository, subcategoryService subcategory.Service) Service {
	return &categoryService{repository, subcategoryService}
}

func (c *categoryService) ListCategory(ctx context.Context) ([]Category, errors.CustomError) {
	return c.repository.ListCategory(ctx)
}

func (c *categoryService) GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError) {
	return c.repository.GetCategoryById(ctx, oid)
}

func (c *categoryService) AddCategory(ctx *fiber.Ctx, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError) {
	cateDTO.Subcategory = utils.RemoveDuplicateObjectIds(cateDTO.Subcategory)
	subcategories, err := c.subcategoryService.FindByIds3(ctx.Context(), cateDTO.Subcategory)
	if err != nil {
		return cateDTO, err
	}

	cateDTO, err = c.repository.AddCategory(ctx.Context(), cateDTO)
	if err != nil {
		return cateDTO, err
	}

	//TODO: tx this
	for _, subcate := range subcategories {
		_, err = c.subcategoryService.InsertToCategory(ctx, &subcate, cateDTO.ID)
		if err != nil {
			return cateDTO, err
		}
	}
	return cateDTO, err
}

func contains(s []primitive.ObjectID, e primitive.ObjectID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (c *categoryService) EditCategory(ctx *fiber.Ctx, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError) {
	oldCate, err := c.GetCategoryById(ctx.Context(), cateDTO.ID)
	if err != nil {
		return cateDTO, err
	}

	cateDTO.Subcategory = utils.RemoveDuplicateObjectIds(cateDTO.Subcategory)
	// check if parse subcategoryIds exist
	subcategories, err := c.subcategoryService.FindByIds3(ctx.Context(), cateDTO.Subcategory)
	if err != nil {
		return cateDTO, err
	}

	// find the remove subcategories by checking
	// the non-intersecting of oldSubcategories and updatedSubcategories
	removeSubcategories := []subcategory.Subcategory{}
	for _, oldSubcategory := range oldCate.Subcategory {
		if !contains(cateDTO.Subcategory, oldSubcategory.ID) {
			removeSubcategories = append(removeSubcategories, oldSubcategory)
		}
	}

	cateDTO, err = c.repository.EditCategory(ctx.Context(), cateDTO)
	if err != nil {
		return cateDTO, err
	}

	// set categoryId to the updated ones
	//TODO: tx this
	for _, subcate := range subcategories {
		_, err = c.subcategoryService.InsertToCategory(ctx, &subcate, cateDTO.ID)
		if err != nil {
			return cateDTO, err
		}
	}

	// unset the remove ones
	for _, removeSubcate := range removeSubcategories {
		_, err = c.subcategoryService.InsertToCategory(ctx, &removeSubcate, primitive.NilObjectID)
		if err != nil {
			return cateDTO, err
		}
	}

	return cateDTO, nil

}

func (c *categoryService) DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	return c.repository.DeleteCategory(ctx, oid)
}

func (c *categoryService) SearchCategory(ctx *fiber.Ctx) ([]CategorySearch, errors.CustomError) {
	return c.repository.SearchCategory(ctx.Context())
}

func (c *categoryService) AddSubcategory(ctx *fiber.Ctx, oid primitive.ObjectID, sid primitive.ObjectID) (*CategoryDTO, errors.CustomError) {
	cateDTO, err := c.repository.GetCategoryByIdNoLookup(ctx.Context(), oid)
	if err != nil {
		return cateDTO, err
	}

	// find new added sid
	addedSubcate, err := c.subcategoryService.GetSubcategoryById(ctx.Context(), sid)
	if err != nil {
		return cateDTO, err
	}

	// append new one to the existing one
	cateDTO.Subcategory = append(cateDTO.Subcategory, sid)
	cateDTO.Subcategory = utils.RemoveDuplicateObjectIds(cateDTO.Subcategory)
	cateDTO, err = c.repository.EditCategory(ctx.Context(), cateDTO)
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
