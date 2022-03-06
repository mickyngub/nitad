package search

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	SearchAll(ctx *fiber.Ctx) (Search, errors.CustomError)
}

type searchService struct {
	categoryService category.Service
	projectService  project.Service
}

func NewService(categoryService category.Service, projectService project.Service) Service {
	return &searchService{
		categoryService,
		projectService,
	}
}

func (s *searchService) SearchAll(ctx *fiber.Ctx) (Search, errors.CustomError) {
	search := Search{}

	categorySearch, err := s.categoryService.SearchCategory(ctx)
	if err != nil {
		return search, err
	}

	projectSearch, err := s.projectService.SearchProject(ctx)
	if err != nil {
		return search, err
	}

	search.Category = categorySearch
	search.Project = projectSearch

	return search, nil
}
