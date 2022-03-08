package project

import (
	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	projectRoute fiber.Router,
) {

	controller := &Controller{service}

	projectRoute.Get("/", controller.ListProject)
	projectRoute.Get("/:projectId", controller.GetProjectById)

	projectRoute.Use(admin.IsAuth())
	projectRoute.Post("/", AddAndEditProjectValidator, controller.AddProject)
	projectRoute.Put("/:projectId", AddAndEditProjectValidator, controller.EditProject)
	projectRoute.Delete("/:projectId", controller.DeleteProject)

}

type Controller struct {
	service Service
}

// list all projects
func (c *Controller) ListProject(ctx *fiber.Ctx) error {
	pq := new(ProjectQuery)
	if err := ctx.QueryParser(pq); err != nil {
		return err
	}

	projects, paginate, err := c.service.ListProject(ctx, pq)
	if err != nil {
		return errors.Throw(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": projects, "paginate": paginate})
}

// get project by id
func (c *Controller) GetProjectById(ctx *fiber.Ctx) error {
	projectId := ctx.Params("projectId")

	project, err := c.service.GetProjectById(ctx, projectId)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": project})
}

// add a project
func (c *Controller) AddProject(ctx *fiber.Ctx) error {
	projectDTO := new(ProjectDTO)
	ctx.BodyParser(projectDTO)

	reports, err := utils.ExtractFiles(ctx, "report")
	if err != nil {
		return errors.Throw(ctx, err)
	}

	images, err := utils.ExtractFiles(ctx, "images")
	if err != nil {
		return errors.Throw(ctx, err)
	}

	projectDTO.Images = images
	projectDTO.Report = reports[0]

	addedProject, err := c.service.AddProject(ctx, projectDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": addedProject})
}

func (c *Controller) EditProject(ctx *fiber.Ctx) error {
	projectId := ctx.Params("projectId")

	projectDTO := new(ProjectDTO)
	ctx.BodyParser(projectDTO)

	reports, err := utils.ExtractUpdatedFiles(ctx, "report")
	if err != nil {
		return errors.Throw(ctx, err)
	}

	images, err := utils.ExtractUpdatedFiles(ctx, "images")
	if err != nil {
		return errors.Throw(ctx, err)
	}

	projectDTO.Images = images
	if reports == nil {
		projectDTO.Report = nil
	} else {
		projectDTO.Report = reports[0]
	}

	editedProject, err := c.service.EditProject(ctx, projectId, projectDTO)
	if err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": editedProject})
}

// delete the project
func (c *Controller) DeleteProject(ctx *fiber.Ctx) error {
	projectId := ctx.Params("projectId")

	if err := c.service.DeleteProject(ctx, projectId); err != nil {
		return errors.Throw(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete project " + projectId + " successfully!"})

}
