package project

import (
	"context"
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
)

func (p *projectService) HandleSubcateAndCateConnection(ctx context.Context, projectDTO *ProjectDTO) ([]category.Category, errors.CustomError) {
	_, sids, err := p.subcategoryService.FindByIds(ctx, projectDTO.Subcategory)
	if err != nil {
		return []category.Category{}, err
	}

	categories, _, err := p.categoryService.FindByIds(ctx, projectDTO.Category)
	if err != nil {
		return []category.Category{}, err
	}

	finalCategories, err := p.categoryService.FilterCatesWithSids(categories, sids)
	if err != nil {
		return []category.Category{}, err
	}
	return finalCategories, nil
}

func (p *projectService) HandleUpdateImages(ctx context.Context, oldImageFilenames []string, newUploadImages []*multipart.FileHeader, deleteImages []string) ([]string, errors.CustomError) {
	imageFilenames := oldImageFilenames

	// DELETE IMAGES
	if len(deleteImages) > 0 {
		deleteFilenames := []string{}
		for _, deleteImage := range deleteImages {
			deleteFilename := gcp.GetFilepath(deleteImage)
			deleteFilenames = append(deleteFilenames, deleteFilename)
		}
		imageFilenames = utils.RemoveSliceFromSlice(imageFilenames, deleteFilenames)
		p.gcpService.DeleteFiles(ctx, deleteFilenames)

	}

	// UPLOAD NEW IMAGE FILES
	if len(newUploadImages) > 0 {
		newImageFilenames, err := p.gcpService.UploadFiles(ctx, newUploadImages, collectionName)
		if err != nil {
			p.gcpService.DeleteFiles(ctx, newImageFilenames)
			return imageFilenames, err
		}
		imageFilenames = append(imageFilenames, newImageFilenames...)
	}

	return imageFilenames, nil
}

func (p *projectService) HandleUpdateReport(ctx context.Context, oldReportURL string, newReportFile *multipart.FileHeader) (string, errors.CustomError) {
	if newReportFile == nil {
		return oldReportURL, nil
	}

	p.gcpService.DeleteFile(ctx, oldReportURL)
	newUploadReportURL, err := p.gcpService.UploadFile(ctx, newReportFile, collectionName)
	if err != nil {
		p.gcpService.DeleteFile(ctx, newUploadReportURL)
		return oldReportURL, err
	}
	return newUploadReportURL, nil
}

func (p *projectService) GetAllURLs(project *Project) {
	images := []string{}
	for _, image := range project.Images {
		images = append(images, gcp.GetURL(image))
	}
	project.Images = images
	project.Report = gcp.GetURL(project.Report)

	for _, cate := range project.Category {
		for _, subcate := range cate.Subcategory {
			subcate.Image = gcp.GetURL(subcate.Image)
		}
	}
}
