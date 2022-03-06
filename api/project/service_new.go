package project

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
)

type Service interface {
	ListProject(ctx context.Context, pq *ProjectQuery) ([]Project, *paginate.Paginate, errors.CustomError)
}

type projectService struct {
	repository         Repository
	subcategoryService subcategory.Service
	categoryService    category.Service
	gcpService         gcp.Uploader
}

func NewService(repository Repository, subcategoryService subcategory.Service, categoryService category.Service, gcpService gcp.Uploader) Service {
	return &projectService{
		repository,
		subcategoryService,
		categoryService,
		gcpService,
	}
}

func (p *projectService) ListProject(ctx context.Context, pq *ProjectQuery) ([]Project, *paginate.Paginate, errors.CustomError) {
	_, sids, err := p.subcategoryService.FindByIds2(ctx, pq.SubcategoryId)
	if err != nil {
		return []Project{}, nil, err
	}

	return p.repository.ListProject(ctx, pq, sids)
}
