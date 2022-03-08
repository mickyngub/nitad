package project

import (
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func GetProjectValidator(c *fiber.Ctx) error {
	projectId := c.Params("projectId")

	HandleCacheGetProjectById(c, projectId)

	return c.Next()
}

func AddAndEditProjectValidator(ctx *fiber.Ctx) error {
	projectDTO := new(ProjectDTO)

	if err := ctx.BodyParser(projectDTO); err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(projectDTO)
	if err != nil {
		return errors.Throw(ctx, errors.NewBadRequestError(err.Error()))
	}

	return ctx.Next()
}
