package project

import (
	"os"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListProject(ctx *fiber.Ctx, pq *ProjectQuery) ([]Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx *fiber.Ctx, oid primitive.ObjectID) (*Project, errors.CustomError)
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

func (p *projectService) ListProject(ctx *fiber.Ctx, pq *ProjectQuery) ([]Project, *paginate.Paginate, errors.CustomError) {
	_, sids, err := p.subcategoryService.FindByIds2(ctx.Context(), pq.SubcategoryId)
	if err != nil {
		return []Project{}, nil, err
	}

	return p.repository.ListProject(ctx.Context(), pq, sids)
}

func (p *projectService) GetProjectById(ctx *fiber.Ctx, oid primitive.ObjectID) (*Project, errors.CustomError) {
	if os.Getenv("APP_ENV") != "test" {
		cacheProject := HandleCacheGetProjectById(ctx, oid)
		if cacheProject != nil {
			return cacheProject, nil
		}
	}

	project, err := p.repository.GetProjectById(ctx.Context(), oid)
	if err != nil {
		return nil, err
	}

	p.repository.IncrementView(ctx.Context(), oid, 1)

	if os.Getenv("APP_ENV") != "test" {
		redis.SetCache(ctx.Path(), project)
	}

	return project, nil
}
