package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
)

type Service interface {
	ListCategory(ctx context.Context) ([]Category, errors.CustomError)
	// GetCateListCategoryById(ctx context.Context, oid primitive.ObjectID) (*CateListCategory, errors.CustomError)
	// AddCateListCategory(ctx context.Context, files []*multipart.FileHeader, subcate *CateListCategory) (*CateListCategory, errors.CustomError)
	// EditCateListCategory(ctx *fiber.Ctx, subcate *CateListCategory) (*CateListCategory, errors.CustomError)
	// DeleteCateListCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError
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
