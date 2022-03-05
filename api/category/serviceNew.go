package category

import (
	"context"
	"time"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListCategory(ctx context.Context) ([]Category, errors.CustomError)
	GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	AddCategory(ctx context.Context, cateDTO *CategoryDTO) (*Category, errors.CustomError)
	EditCategory(ctx context.Context, cateDTO *CategoryDTO) (*Category, errors.CustomError)
	// DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError
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

func (c *categoryService) AddCategory(ctx context.Context, cateDTO *CategoryDTO) (*Category, errors.CustomError) {
	var cate Category
	subcategories, sids, err := c.subcategoryService.FindByIds2(ctx, cateDTO.Subcategory)
	if err != nil {
		return &cate, err
	}

	utils.CopyStruct(cateDTO, cate)

	now := time.Now()
	cate.CreatedAt = now
	cate.UpdatedAt = now
	cate.Subcategory = subcategories
	return c.repository.AddCategory(ctx, sids, &cate)
}

func (c *categoryService) EditCategory(ctx context.Context, cateDTO *CategoryDTO) (*Category, errors.CustomError) {
	var cate Category
	subcategories, sids, err := c.subcategoryService.FindByIds2(ctx, cateDTO.Subcategory)
	if err != nil {
		return &cate, err
	}

	utils.CopyStruct(cateDTO, cate)

	now := time.Now()
	cate.UpdatedAt = now
	cate.Subcategory = subcategories
	return c.repository.AddCategory(ctx, sids, &cate)
}
