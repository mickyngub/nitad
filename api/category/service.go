package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListCategory(ctx context.Context) ([]Category, errors.CustomError)
	GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	AddCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError)
	EditCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError)
	DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	FindByIds2(ctx context.Context, cids []string) ([]Category, []primitive.ObjectID, errors.CustomError)
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

func (c *categoryService) AddCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError) {
	_, osids, err := c.subcategoryService.FindByIds2(ctx, cateDTO.Subcategory)
	if err != nil {
		return cateDTO, err
	}

	return c.repository.AddCategory(ctx, cateDTO, osids)
}

func (c *categoryService) EditCategory(ctx context.Context, cateDTO *CategoryDTO) (*CategoryDTO, errors.CustomError) {
	_, osids, err := c.subcategoryService.FindByIds2(ctx, cateDTO.Subcategory)
	if err != nil {
		return cateDTO, err
	}

	return c.repository.EditCategory(ctx, cateDTO, osids)
}

func (c *categoryService) DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	return c.repository.DeleteCategory(ctx, oid)
}
