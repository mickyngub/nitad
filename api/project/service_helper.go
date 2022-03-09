package project

import (
	"log"
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (p *projectService) HandleSubcateAndCateConnection(ctx *fiber.Ctx, projectDTO *ProjectDTO) ([]category.Category, errors.CustomError) {
	_, soids, err := p.subcategoryService.FindByIds3(ctx.Context(), projectDTO.Subcategory)
	if err != nil {
		return []category.Category{}, err
	}

	categories, _, err := p.categoryService.FindByIds2(ctx.Context(), projectDTO.Category)
	if err != nil {
		return []category.Category{}, err
	}

	finalCategories, err := category.FilterCatesWithSids(categories, soids)
	if err != nil {
		return []category.Category{}, err
	}
	return finalCategories, nil
}

func (p *projectService) HandleUpdateImages(ctx *fiber.Ctx, oldImageFilenames []string, newUploadImages []*multipart.FileHeader, deleteImages []string) ([]string, errors.CustomError) {
	imageFilenames := oldImageFilenames

	// DELETE IMAGES
	if len(deleteImages) > 0 {
		deleteFilenames := []string{}
		for _, deleteImage := range deleteImages {
			deleteFilename := gcp.GetFilepath(deleteImage)
			zap.S().Info("Deleted", deleteFilename)

			deleteFilenames = append(deleteFilenames, deleteFilename)
		}
		imageFilenames = utils.RemoveSliceFromSlice(imageFilenames, deleteFilenames)
		zap.S().Info("Deleted", imageFilenames)
		p.gcpService.DeleteFiles(ctx.Context(), deleteFilenames)
	}

	// UPLOAD NEW IMAGE FILES
	if len(newUploadImages) > 0 {
		newImageFilenames, err := p.gcpService.UploadFiles(ctx.Context(), newUploadImages, collectionName)
		zap.S().Info("pass 6", newImageFilenames)
		if err != nil {
			p.gcpService.DeleteFiles(ctx.Context(), newImageFilenames)
			return imageFilenames, err
		}
		imageFilenames = append(imageFilenames, newImageFilenames...)
	}

	return imageFilenames, nil
}

func (p *projectService) HandleUpdateReport(ctx *fiber.Ctx, oldReportURL string, newReportFile *multipart.FileHeader) (string, errors.CustomError) {
	if newReportFile == nil {
		return oldReportURL, nil
	}

	p.gcpService.DeleteFile(ctx.Context(), oldReportURL)
	newUploadReportURL, err := p.gcpService.UploadFile(ctx.Context(), newReportFile, collectionName)
	if err != nil {
		p.gcpService.DeleteFile(ctx.Context(), newUploadReportURL)
		return oldReportURL, err
	}
	return newUploadReportURL, nil
}

func (p *projectService) GetAllURLs(project *Project) {
	images := []string{}
	for _, image := range project.Images {
		log.Println("ss", gcp.GetURL(image))
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
