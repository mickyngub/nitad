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
	SearchProject(ctx *fiber.Ctx) ([]ProjectSearch, errors.CustomError)
	ListProject(ctx *fiber.Ctx, pq *ProjectQuery) ([]Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx *fiber.Ctx, oid primitive.ObjectID) (*Project, errors.CustomError)
	AddProject(ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	EditProject(ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	DeleteProject(ctx *fiber.Ctx, oid primitive.ObjectID) errors.CustomError
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

func (p *projectService) SearchProject(ctx *fiber.Ctx) ([]ProjectSearch, errors.CustomError) {
	return p.repository.SearchProject(ctx.Context())
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

	finalCategories, err := p.HandleSubcateAndCateConnection(ctx, projectDTO)
	if err != nil {
		return project, err
	}

	reportURL, err := p.gcpService.UploadFile(ctx.Context(), projectDTO.Report, collectionName)
	if err != nil {
		return project, err
	}

	imageURLs, err := p.gcpService.UploadFiles(ctx.Context(), projectDTO.Images, collectionName)
	if err != nil {
		return project, err
	}

	err = utils.CopyStruct(projectDTO, project)
	if err != nil {
		return project, err
	}
	project.Images = imageURLs
	project.Report = reportURL
	project.Category = finalCategories

	addedProject, err := p.repository.AddProject(ctx.Context(), project)
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx.Context(), URLs, collectionName)
		return project, err
	}

	return addedProject, nil
}

func (p *projectService) EditProject(ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError) {
	editedProject := new(Project)

	finalCategories, err := p.HandleSubcateAndCateConnection(ctx, projectDTO)
	if err != nil {
		return editedProject, err
	}

	oldProj, err := p.GetProjectById(ctx, projectDTO.ID)
	if err != nil {
		return editedProject, err
	}

	imageURLs, err := p.HandleUpdateImages(ctx, oldProj.Images, projectDTO.Images, projectDTO.DeleteImages)
	if err != nil {
		return editedProject, err
	}

	reportURL, err := p.HandleUpdateReport(ctx, oldProj.Report, projectDTO.Report)
	if err != nil {
		return editedProject, err
	}

	err = utils.CopyStruct(projectDTO, editedProject)
	if err != nil {
		return editedProject, err
	}
	editedProject.Images = imageURLs
	editedProject.Report = reportURL
	editedProject.Category = finalCategories

	editedProject, err = p.repository.EditProject(ctx.Context(), editedProject)
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx.Context(), URLs, collectionName)
		return editedProject, err
	}
	return editedProject, nil
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
