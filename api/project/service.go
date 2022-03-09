package project

import (
	"fmt"
	"os"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	SearchProject(ctx *fiber.Ctx) ([]ProjectSearch, errors.CustomError)
	ListProject(ctx *fiber.Ctx, pq *ProjectQuery) ([]*Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx *fiber.Ctx, id string) (*Project, errors.CustomError)
	AddProject(ctx *fiber.Ctx, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	EditProject(ctx *fiber.Ctx, id string, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	DeleteProject(ctx *fiber.Ctx, id string) errors.CustomError

	GetAllURLs(project *Project)
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

//TODO: just ignore findbyIds subcate
// if not found then nothing happen
func (p *projectService) ListProject(ctx *fiber.Ctx, pq *ProjectQuery) ([]*Project, *paginate.Paginate, errors.CustomError) {
	_, sids, err := p.subcategoryService.FindByIds3(ctx.Context(), pq.SubcategoryId)
	if err != nil {
		return nil, nil, err
	}

	projects, pagin, err := p.repository.ListProject(ctx.Context(), pq, sids)
	if err != nil {
		return nil, nil, err
	}

	for _, project := range projects {
		p.GetAllURLs(project)
	}

	return projects, pagin, nil
}

func (p *projectService) GetProjectById(ctx *fiber.Ctx, id string) (*Project, errors.CustomError) {
	if os.Getenv("APP_ENV") != "test" {
		cacheProject := HandleCacheGetProjectById(ctx, id)
		if cacheProject != nil {
			return cacheProject, nil
		}
	}

	project, err := p.repository.GetProjectById(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	p.repository.IncrementView(ctx.Context(), project.ID, 1)

	p.GetAllURLs(project)

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
	fmt.Println("add", project.Category[0].Subcategory[0].Image)

	addedProject, err := p.repository.AddProject(ctx.Context(), project)
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx.Context(), URLs)
		return project, err
	}

	return addedProject, nil
}

func (p *projectService) EditProject(ctx *fiber.Ctx, id string, projectDTO *ProjectDTO) (*Project, errors.CustomError) {
	editedProject := new(Project)

	finalCategories, err := p.HandleSubcateAndCateConnection(ctx, projectDTO)
	if err != nil {
		return editedProject, err
	}

	oldProj, err := p.repository.GetProjectById(ctx.Context(), id)
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

	editedProject.ID = oldProj.ID
	editedProject.Images = imageURLs
	editedProject.Report = reportURL
	editedProject.Category = finalCategories

	editedProject, err = p.repository.EditProject(ctx.Context(), editedProject)
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx.Context(), URLs)
		return editedProject, err
	}
	return editedProject, nil
}

func (p *projectService) DeleteProject(ctx *fiber.Ctx, id string) errors.CustomError {
	project, err := p.repository.GetProjectById(ctx.Context(), id)
	if err != nil {
		return err
	}

	p.gcpService.DeleteFile(ctx.Context(), project.Report)
	p.gcpService.DeleteFiles(ctx.Context(), project.Images)

	return p.repository.DeleteProject(ctx.Context(), project.ID)
}
