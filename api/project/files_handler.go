package project

import (
	"context"
	"log"

	"github.com/birdglove2/nitad-backend/api/collections_helper"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleUpdateReportAndImages(c *fiber.Ctx, updateProject *Project) (*Project, errors.CustomError) {
	oldProject, err := GetById(updateProject.ID)
	if err != nil {
		return updateProject, err
	}

	// assign attrs from DB to the updateOne in case there is no update
	updateProject.Report = oldProject.Report
	updateProject.Images = oldProject.Images

	log.Println("Report", updateProject.Report)

	updateProject, err = HandleUpdateImages(c, updateProject)
	if err != nil {
		return updateProject, err
	}

	updateProject, err = HandleUpdateReport(c, updateProject)
	if err != nil {
		return updateProject, err
	}

	updateProject.CreatedAt = oldProject.CreatedAt
	return updateProject, nil
}

func HandleUpdateReport(c *fiber.Ctx, p *Project) (*Project, errors.CustomError) {
	files, err := utils.ExtractUpdatedFiles(c, "report")
	if err != nil {
		return p, err
	}

	// if there is file passed, delete the old one and upload a new one
	if len(files) > 0 {
		newUploadFilename, err := collections_helper.HandleUpdateSingleFile(c.Context(), files[0], p.Report, collectionName)
		if err != nil {
			return p, err
		}
		// if upload success, pass the url to the project struct
		p.Report = newUploadFilename
	}

	return p, nil
}

func HandleUpdateImages(c *fiber.Ctx, up *Project) (*Project, errors.CustomError) {
	// DELETE FILES
	if len(up.DeleteImages) > 0 {
		// remove deleteImages from Images attrs
		up.Images = utils.RemoveSliceFromSlice(up.Images, up.DeleteImages)
		gcp.DeleteFiles(c.Context(), up.DeleteImages, collectionName)

	}

	// UPLOAD NEW FILES
	files, err := utils.ExtractUpdatedFiles(c, "images")
	if err != nil {
		return up, err
	}
	if len(files) > 0 {
		// if file pass, upload file
		imageURLs, err := gcp.UploadFiles(c.Context(), files, collectionName)

		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			gcp.DeleteFiles(c.Context(), imageURLs, collectionName)
			return up, err
		}

		// concat uploaded file to the existing ones
		up.Images = append(up.Images, imageURLs...)
	}

	return up, nil
}

func HandleDeleteImages(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	project, err := GetById(oid)
	if err != nil {
		return err
	}

	gcp.DeleteFiles(ctx, project.Images, collectionName)

	return nil
}
