package project

import (
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/errors"
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

func (p *projectService) HandleUpdateImages(ctx *fiber.Ctx, oldImageURLs []string, newUploadImages []*multipart.FileHeader, deleteImages []string) ([]string, errors.CustomError) {
	imageURLs := oldImageURLs
	zap.S().Info("pass 5", imageURLs)

	// DELETE IMAGES
	if len(deleteImages) > 0 {
		imageURLs = utils.RemoveSliceFromSlice(imageURLs, deleteImages)
		zap.S().Info("Deleted", imageURLs)
		p.gcpService.DeleteFiles(ctx.Context(), deleteImages, collectionName)
	}

	// UPLOAD NEW IMAGE FILES
	if len(newUploadImages) > 0 {

		newImageURLs, err := p.gcpService.UploadFiles(ctx.Context(), newUploadImages, collectionName)
		zap.S().Info("pass 6", newImageURLs)
		if err != nil {
			p.gcpService.DeleteFiles(ctx.Context(), newImageURLs, collectionName)
			return imageURLs, err
		}
		imageURLs = append(imageURLs, newImageURLs...)
	}

	return imageURLs, nil
}

func (p *projectService) HandleUpdateReport(ctx *fiber.Ctx, oldReportURL string, newReportFile *multipart.FileHeader) (string, errors.CustomError) {
	if newReportFile == nil {
		return oldReportURL, nil
	}

	p.gcpService.DeleteFile(ctx.Context(), oldReportURL, collectionName)
	newUploadReportURL, err := p.gcpService.UploadFile(ctx.Context(), newReportFile, collectionName)
	if err != nil {
		p.gcpService.DeleteFile(ctx.Context(), newUploadReportURL, collectionName)
		return oldReportURL, err
	}
	return newUploadReportURL, nil
}
