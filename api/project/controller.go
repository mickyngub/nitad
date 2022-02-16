package project

import (
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewController(
	projectRoute fiber.Router,
) {

	controller := &Controller{}

	projectRoute.Get("/", controller.ListProject)
	projectRoute.Get("/:projectId", controller.GetProject)

	//TODO Add auth
	projectRoute.Post("/", controller.AddProject)
	projectRoute.Put("/:projectId", controller.EditProject)
	projectRoute.Delete("/:projectId", controller.DeleteProject)

}

type Controller struct{}

var collectionName = database.COLLECTIONS["PROJECT"]

type ProjectQuery struct {
	SubcategoryId []string `query:"subcategoryId`
}

// list all projects
// TODO: sorting & filtering
func (contc *Controller) ListProject(c *fiber.Ctx) error {
	p := new(ProjectQuery)

	if err := c.QueryParser(p); err != nil {
		return err
	}

	validSubcategoryIds, err := subcategory.ValidateIds(p.SubcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	projects, err := FindAll(validSubcategoryIds)
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
	p := new(ProjectRequest)
	if err := c.BodyParser(p); err != nil {
		return errors.Throw(c, errors.InvalidInput)
	}

	files, err := functions.ExtractFiles(c, "images")
	if err != nil {
		return errors.Throw(c, err)
	}

	imageURLs, err := gcp.UploadImages(files, collectionName)
	if err != nil {
		return errors.Throw(c, err)
	}
	p.Images = imageURLs

	result, err := Add(p)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

func (contc *Controller) EditProject(c *fiber.Ctx) error {
	upr := new(UpdateProjectRequest)

	if err := c.BodyParser(upr); err != nil {
		return errors.Throw(c, errors.InvalidInput)
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
