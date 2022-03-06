package project

import (
	"os"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListProject(ctx *fiber.Ctx, pq *ProjectQuery) ([]Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx *fiber.Ctx, oid primitive.ObjectID) (*Project, errors.CustomError)
	AddProject(ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	EditProject(c *Controller, ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	DeleteProject(ctx *fiber.Ctx, oid primitive.ObjectID) errors.CustomError
	SearchProject(ctx *fiber.Ctx) ([]ProjectSearch, errors.CustomError)

	// HandleUpdateReportAndImages(ctx *fiber.Ctx, proj *Project) (*Project, errors.CustomError)
	// HandleUpdateImages(ctx *fiber.Ctx, proj *Project) (*Project, errors.CustomError)
	// HandleUpdateReport(ctx *fiber.Ctx, proj *Project) (*Project, errors.CustomError)
	// HandleDeleteImages(ctx *fiber.Ctx, oid primitive.ObjectID) errors.CustomError
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

func (p *projectService) AddProject(ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError) {
	project := new(Project)

	_, sids, err := p.subcategoryService.FindByIds2(ctx.Context(), projectDTO.Subcategory)
	if err != nil {
		return project, err
	}

	categories, _, err := p.categoryService.FindByIds2(ctx.Context(), projectDTO.Category)
	if err != nil {
		return project, err
	}

	finalCategories, err := category.FilterCatesWithSids(categories, sids)
	if err != nil {
		return project, err
	}

	files, err := utils.ExtractFiles(ctx, "report")
	if err != nil {
		return project, err
	}

	reportURL, err := p.gcpService.UploadFile(ctx.Context(), files[0], collectionName)
	if err != nil {
		return project, err
	}

	files, err = utils.ExtractFiles(ctx, "images")
	if err != nil {
		return project, err
	}
	imageURLs, err := p.gcpService.UploadFiles(ctx.Context(), files, collectionName)
	if err != nil {
		return project, err
	}

	err = utils.CopyStruct(projectDTO, project)
	if err != nil {
		return project, err
	}

	project.Category = finalCategories
	project.Images = imageURLs
	project.Report = reportURL

	addedProject, err := p.repository.AddProject(ctx.Context(), project)
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx.Context(), URLs, collectionName)
		return project, err
	}

	return addedProject, nil
}

func (p *projectService) EditProject(c *Controller, ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError) {
	project := new(Project)
	err := utils.CopyStruct(projectDTO, project)
	if err != nil {
		return project, err
	}

	project, err = p.HandleUpdateReportAndImages(ctx, project)
	if err != nil {
		return project, err
	}

	return p.repository.EditProject(ctx.Context(), project)
}

func (p *projectService) DeleteProject(ctx *fiber.Ctx, oid primitive.ObjectID) errors.CustomError {
	project, err := p.GetProjectById(ctx, oid)
	if err != nil {
		return err
	}

	p.gcpService.DeleteFile(ctx.Context(), project.Report, collectionName)
	p.gcpService.DeleteFiles(ctx.Context(), project.Images, collectionName)

	return p.repository.DeleteProject(ctx.Context(), oid)

}

func (p *projectService) SearchProject(ctx *fiber.Ctx) ([]ProjectSearch, errors.CustomError) {
	return p.repository.SearchProject(ctx.Context())
}
