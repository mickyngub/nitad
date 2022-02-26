package project

import (
	"strconv"

	"github.com/birdglove2/nitad-backend/api/admin"
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
	projectRoute.Post("/", AddProjectValidator, controller.AddProject)
	projectRoute.Put("/:projectId", EditProjectValidator, controller.EditProject)
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

	queryString := pq.Sort + strconv.Itoa(pq.By) + strconv.Itoa(pq.Page) + strconv.Itoa(pq.Limit)

	for _, sid := range pq.SubcategoryId {
		queryString += sid
	}

	var p []*Project
	redis.GetCache(queryString, &p)

	if p != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": p})
	}

	projects, paginate, err := FindAll(pq)
	if err != nil {
		return errors.Throw(c, err)
	}

	redis.SetCache(queryString, projects)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": projects, "paginate": paginate})
}

// get project by id
func (contc *Controller) GetProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")

	objectId, err := utils.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var p *Project
	redis.GetCache(projectId, &p)

	if p != nil {
		IncrementViewCache(projectId, p.Views)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": p})
	}

	var result Project
	if result, err = GetById(objectId); err != nil {
		return errors.Throw(c, err)
	}

	redis.SetCache(projectId, result)

	IncrementView(objectId, 1)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a project
func (contc *Controller) AddProject(c *fiber.Ctx) error {
	projectBody, ok := c.Locals("projectBody").(*Project)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Add project went wrong!"))
	}

	files, err := utils.ExtractFiles(c, "images")
	if err != nil {
		return errors.Throw(c, err)
	}

	imageURLs, err := gcp.UploadImages(c.Context(), files, collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}
	projectBody.Images = imageURLs

	result, err := Add(projectBody)
	if err != nil {
		// if there is any error, remove the uploaded file from gcp
		gcp.DeleteImages(c.Context(), imageURLs, collectionName)
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

func (contc *Controller) EditProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")
	projectIdObjectId, err := utils.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	updateProjectBody, ok := c.Locals("updateProjectBody").(*UpdateProject)
	if !ok {
		return errors.Throw(c, errors.NewInternalServerError("Edit project went wrong!"))
	}

	updateProjectBody, err = HandleUpdateImages(c, updateProjectBody, projectIdObjectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	result, err := Edit(projectIdObjectId, updateProjectBody)
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

	err = HandleDeleteImages(c.Context(), objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	err = Delete(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete project " + projectId + " successfully!"})

}
