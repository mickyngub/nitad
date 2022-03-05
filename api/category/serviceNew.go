package category

import "github.com/birdglove2/nitad-backend/api/subcategory"

type Service interface{}

type categoryService struct {
	repository         Repository
	subcategoryService subcategory.Service
}

func NewService(repository Repository, subcategoryService subcategory.Service) Service {
	return &categoryService{repository, subcategoryService}
}
