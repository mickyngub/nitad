package project

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/collections_helper"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (contc *Controller) HandleUpdateReportAndImages(c *fiber.Ctx, updateProject *Project) (*Project, errors.CustomError) {
	oldProject, err := GetById(updateProject.ID)
	if err != nil {
		return updateProject, err
	}

	// assign attrs from DB to the updateOne in case there is no update
	updateProject.Report = oldProject.Report
	updateProject.Images = oldProject.Images

	updateProject, err = contc.HandleUpdateImages(c, updateProject)
	if err != nil {
		return updateProject, err
	}

	updateProject, err = contc.HandleUpdateReport(c, updateProject)
	if err != nil {
		return updateProject, err
	}

	updateProject.CreatedAt = oldProject.CreatedAt
	return updateProject, nil
}

func (contc *Controller) HandleUpdateReport(c *fiber.Ctx, p *Project) (*Project, errors.CustomError) {
	files, err := utils.ExtractUpdatedFiles(c, "report")
	if err != nil {
		return p, err
	}

	// if there is file passed, delete the old one and upload a new one
	if len(files) > 0 {
		newUploadFilename, err := collections_helper.HandleUpdateSingleFile(contc.gcpService, c.Context(), files[0], p.Report, collectionName)
		if err != nil {
			return p, err
		}
		// if upload success, pass the url to the project struct
		p.Report = newUploadFilename
	}

	return p, nil
}

func (contc *Controller) HandleUpdateImages(c *fiber.Ctx, up *Project) (*Project, errors.CustomError) {
	// DELETE FILES
	if len(up.DeleteImages) > 0 {
		// remove deleteImages from Images attrs
		up.Images = utils.RemoveSliceFromSlice(up.Images, up.DeleteImages)
		contc.gcpService.DeleteFiles(c.Context(), up.DeleteImages, collectionName)

	}

	// UPLOAD NEW FILES
	files, err := utils.ExtractUpdatedFiles(c, "images")
	if err != nil {
		return up, err
	}
	if len(files) > 0 {
		// if file pass, upload file
		imageURLs, err := contc.gcpService.UploadFiles(c.Context(), files, collectionName)

		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			contc.gcpService.DeleteFiles(c.Context(), imageURLs, collectionName)
			return up, err
		}

		// concat uploaded file to the existing ones
		up.Images = append(up.Images, imageURLs...)
	}

	return up, nil
}

func (contc *Controller) HandleDeleteImages(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	project, err := GetById(oid)
	if err != nil {
		return err
	}

	contc.gcpService.DeleteFiles(ctx, project.Images, collectionName)

	return nil
}
