package project

import (
	"github.com/birdglove2/nitad-backend/api/collections_helper"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p *projectService) HandleUpdateReportAndImages(ctx *fiber.Ctx, proj *Project) (*Project, errors.CustomError) {
	oldProj, err := p.GetProjectById(ctx, proj.ID)
	if err != nil {
		return proj, err
	}

	proj.Report = oldProj.Report
	proj.Images = oldProj.Images
	proj.CreatedAt = oldProj.CreatedAt

	proj, err = p.HandleUpdateImages(ctx, proj)
	if err != nil {
		return proj, err
	}

	proj, err = p.HandleUpdateReport(ctx, proj)
	if err != nil {
		return proj, err
	}

	return proj, nil
}

func (p *projectService) HandleUpdateImages(ctx *fiber.Ctx, proj *Project) (*Project, errors.CustomError) {
	// DELETE FILES
	if len(proj.DeleteImages) > 0 {
		// remove deleteImages from Images attrs
		proj.Images = utils.RemoveSliceFromSlice(proj.Images, proj.DeleteImages)
		p.gcpService.DeleteFiles(ctx.Context(), proj.DeleteImages, collectionName)
	}

	// UPLOAD NEW FILES
	files, err := utils.ExtractUpdatedFiles(ctx, "images")
	if err != nil {
		return proj, err
	}
	if len(files) > 0 {
		// if file pass, upload file
		imageURLs, err := p.gcpService.UploadFiles(ctx.Context(), files, collectionName)

		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			p.gcpService.DeleteFiles(ctx.Context(), imageURLs, collectionName)
			return proj, err
		}

		// concat uploaded file to the existing ones
		proj.Images = append(proj.Images, imageURLs...)
	}

	return proj, nil
}

func (p *projectService) HandleUpdateReport(ctx *fiber.Ctx, proj *Project) (*Project, errors.CustomError) {
	files, err := utils.ExtractUpdatedFiles(ctx, "report")
	if err != nil {
		return proj, err
	}

	// if there is file passed, delete the old one and upload a new one
	if len(files) > 0 {
		newUploadFilename, err := collections_helper.HandleUpdateSingleFile(p.gcpService, ctx.Context(), files[0], proj.Report, collectionName)
		if err != nil {
			return proj, err
		}
		// if upload success, pass the url to the project struct
		proj.Report = newUploadFilename
	}

	return proj, nil
}

func (p *projectService) HandleDeleteImages(ctx *fiber.Ctx, oid primitive.ObjectID) errors.CustomError {
	project, err := p.GetProjectById(ctx, oid)
	if err != nil {
		return err
	}

	p.gcpService.DeleteFiles(ctx.Context(), project.Images, collectionName)
	return nil
}
