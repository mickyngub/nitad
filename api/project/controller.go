package project

import (
	"log"

	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	projectRoute fiber.Router,
) {

	controller := &Controller{}

	projectRoute.Get("/", controller.ListProject)
	projectRoute.Get("/:projectId", controller.GetProject)

	projectRoute.Use(admin.IsAuth())
	projectRoute.Post("/", AddAndEditProjectValidator, controller.AddProject)
	projectRoute.Put("/:projectId", AddAndEditProjectValidator, controller.EditProject)
	projectRoute.Delete("/:projectId", controller.DeleteProject)

}

type Controller struct{}

var collectionName = database.COLLECTIONS["PROJECT"]

type ProjectQuery struct {
	SubcategoryId []string `query:"subcategoryId"`
	Sort          string   `query:"sort"`
	By            int      `query:"by"`
	Page          int      `query:"page"`
	Limit         int      `query:"limit"`
}

// list all projects
func (contc *Controller) ListProject(c *fiber.Ctx) error {
	pq := new(ProjectQuery)

	if err := c.QueryParser(pq); err != nil {
		return err
	}

	var projects []Project
	var pagin paginate.Paginate
	var err errors.CustomError
	if pq.Limit == -1 {
		projects, pagin, err = FindAllNoLimit()
	} else {
		projects, pagin, err = FindAll(pq)
	}

	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": projects, "paginate": pagin})

}

// get project by id
func (contc *Controller) GetProject(c *fiber.Ctx) error {
	// GetProjectValidator(c)

	projectId := c.Params("projectId")
	// objectId, _ := primitive.ObjectIDFromHex(projectId)

	objectId, err := utils.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}
	// log.Println("5", projectId)

	cacheProject := HandleCacheGetProjectById(c, projectId)
	// log.Println("6", projectId)
	if cacheProject != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": cacheProject})
	}

	log.Println("7")
	result, err := GetById(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	IncrementView(objectId, 1)
	redis.SetCache(c.Path(), result)
	// log.Println("8")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a project
func (contc *Controller) AddProject(c *fiber.Ctx) error {
	projectBody, ok := c.Locals("projectBody").(*Project)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Add project went wrong!"))
	}

	files, err := utils.ExtractFiles(c, "report")
	if err != nil {
		return errors.Throw(c, err)
	}
	reportURL, err := gcp.UploadFile(c.Context(), files[0], collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}

	files, err = utils.ExtractFiles(c, "images")
	if err != nil {
		return errors.Throw(c, err)
	}
	imageURLs, err := gcp.UploadFiles(c.Context(), files, collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}

	projectBody.Report = reportURL
	projectBody.Images = imageURLs

	result, err := Add(projectBody)
	if err != nil {
		// if there is any error, remove the uploaded files from gcp
		imageURLs = append(imageURLs, reportURL)
		gcp.DeleteFiles(c.Context(), imageURLs, collectionName)
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

func (contc *Controller) EditProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")
	oid, err := utils.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	updateProject, ok := c.Locals("projectBody").(*Project)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Edit project went wrong!"))
	}

	updateProject.ID = oid
	updateProject, err = HandleUpdateReportAndImages(c, updateProject)
	if err != nil {
		return errors.Throw(c, err)
	}

	result, err := Edit(updateProject)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// delete the project
func (cont *Controller) DeleteProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")
	objectId, err := utils.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	project, err := GetById(objectId)
	if err != nil {
		return err
	}

	gcp.DeleteFile(c.Context(), project.Report, collectionName)
	gcp.DeleteFiles(c.Context(), project.Images, collectionName)

	err = Delete(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete project " + project.Title + " successfully!"})

}
