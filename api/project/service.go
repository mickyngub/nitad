package project

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
)

type Service interface {
	SearchProject(ctx context.Context) ([]ProjectSearch, errors.CustomError)
	ListProject(ctx context.Context, pq *ProjectQuery) ([]*Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx context.Context, id string) (*Project, errors.CustomError)
	AddProject(ctx context.Context, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	EditProject(ctx context.Context, id string, projectDTO *ProjectDTO) (*Project, errors.CustomError)
	DeleteProject(ctx context.Context, id string) errors.CustomError

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

func (p *projectService) SearchProject(ctx context.Context) ([]ProjectSearch, errors.CustomError) {
	return p.repository.SearchProject(ctx)
}

//TODO: just ignore findbyIds subcate
// if not found then nothing happen
func (p *projectService) ListProject(ctx context.Context, pq *ProjectQuery) ([]*Project, *paginate.Paginate, errors.CustomError) {
	_, sids, err := p.subcategoryService.FindByIds3(ctx, pq.SubcategoryId)
	if err != nil {
		return nil, nil, err
	}

	projects, pagin, err := p.repository.ListProject(ctx, pq, sids)
	if err != nil {
		return nil, nil, err
	}

	for _, project := range projects {
		p.GetAllURLs(project)
	}

	return projects, pagin, nil
}

func (p *projectService) GetProjectById(ctx context.Context, id string) (*Project, errors.CustomError) {

	project, err := p.repository.GetProjectById(ctx, id)
	if err != nil {
		return nil, err
	}

	p.repository.IncrementView(ctx, project.ID, 1)

	p.GetAllURLs(project)

	return project, nil
}

func (p *projectService) AddProject(ctx context.Context, projectDTO *ProjectDTO) (*Project, errors.CustomError) {
	project := new(Project)

	finalCategories, err := p.HandleSubcateAndCateConnection(ctx, projectDTO)
	if err != nil {
		return project, err
	}

	reportURL, err := p.gcpService.UploadFile(ctx, projectDTO.Report, collectionName)
	if err != nil {
		return project, err
	}

	imageURLs, err := p.gcpService.UploadFiles(ctx, projectDTO.Images, collectionName)
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

	err = database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		var txErr errors.CustomError
		project, txErr = p.addProject(sessionContext, project)
		return txErr
	})
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx, URLs)
		return nil, err
	}
	return project, nil

}

func (p *projectService) addProject(ctx context.Context, project *Project) (*Project, errors.CustomError) {
	for _, cate := range project.Category {
		p.categoryService.IncrementProjectCount(ctx, &cate)
	}

	return p.repository.AddProject(ctx, project)
}

func (p *projectService) EditProject(ctx context.Context, id string, projectDTO *ProjectDTO) (*Project, errors.CustomError) {
	editedProject := new(Project)

	finalCategories, err := p.HandleSubcateAndCateConnection(ctx, projectDTO)
	if err != nil {
		return editedProject, err
	}

	oldProj, err := p.repository.GetProjectById(ctx, id)
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

	err = database.ExecTx(ctx, func(sessionContext context.Context) errors.CustomError {
		var txErr errors.CustomError
		editedProject, txErr = p.editProject(sessionContext, editedProject)
		return txErr
	})

	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		URLs := append(imageURLs, reportURL)
		p.gcpService.DeleteFiles(ctx, URLs)
		return nil, err
	}

	return editedProject, nil
}

func (p *projectService) editProject(ctx context.Context, project *Project) (*Project, errors.CustomError) {
	for _, cate := range project.Category {
		p.categoryService.IncrementProjectCount(ctx, &cate)
	}

	return p.repository.EditProject(ctx, project)
}

func (p *projectService) DeleteProject(ctx context.Context, id string) errors.CustomError {
	project, err := p.repository.GetProjectById(ctx, id)
	if err != nil {
		return err
	}

	p.gcpService.DeleteFile(ctx, project.Report)
	p.gcpService.DeleteFiles(ctx, project.Images)

	return p.repository.DeleteProject(ctx, project.ID)
}
