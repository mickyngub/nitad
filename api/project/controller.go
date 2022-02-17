package project

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewController(
	projectRoute fiber.Router,
) {

	controller := &Controller{}

	projectRoute.Get("/", controller.ListProject)
	projectRoute.Get("/:projectId", controller.GetProject)

	//TODO Add auth
	projectRoute.Post("/", AddProjectValidator, controller.AddProject)
	projectRoute.Put("/:projectId", controller.EditProject)
	projectRoute.Delete("/:projectId", controller.DeleteProject)

}

type Controller struct{}

var collectionName = database.COLLECTIONS["PROJECT"]

type ProjectQuery struct {
	SubcategoryId []string `query:"subcategoryId"`
	// ByViews       int      `query:"byViews"`
	// ByName        int      `query:"byName"`
	// ByCreatedAt   int      `query:"byCreatedAt"`
	// ByUpdatedAt   int      `query:"byUpdatedAt"`
	Sort string `query:"sort`
	By   int    `query:by`
}

// list all projects
func (contc *Controller) ListProject(c *fiber.Ctx) error {
	pq := new(ProjectQuery)

	if err := c.QueryParser(pq); err != nil {
		return err
	}

	projects, err := FindAll(pq)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": projects})
}

// get project by id
func (contc *Controller) GetProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")

	objectId, err := functions.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var result bson.M
	if result, err = FindById(objectId); err != nil {
		return errors.Throw(c, err)
	}

	defer IncrementView(objectId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a project
func (contc *Controller) AddProject(c *fiber.Ctx) error {
	projectBody := c.Locals("projectBody").(*Project)
	cid := c.Locals("cid").(primitive.ObjectID)
	sids := c.Locals("sids").([]primitive.ObjectID)

	files, err := functions.ExtractFiles(c, "images")
	if err != nil {
		return errors.Throw(c, err)
	}

	imageURLs, err := gcp.UploadImages(files, collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}
	projectBody.Images = imageURLs

	result, err := Add(projectBody, cid, sids)
	if err != nil {
		// if there is any error, remove the uploaded file from gcp
		defer gcp.DeleteImages(imageURLs, collectionName)
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

func (contc *Controller) EditProject(c *fiber.Ctx) error {
	upr := new(UpdateProjectRequest)

	if err := c.BodyParser(upr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))

	}

	projectId := c.Params("projectId")
	objectId, err := functions.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	upr, err = HandleUpdateImages(c, upr, objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	result, err := Edit(objectId, upr)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// delete the project
func (cont *Controller) DeleteProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")
	objectId, err := functions.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	err = HandleDeleteImages(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	err = Delete(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete project successfully!"})

}
