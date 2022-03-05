package project

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func GetProjectValidator(c *fiber.Ctx) error {
	projectId := c.Params("projectId")

	_, err := utils.IsValidObjectId(projectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	HandleCacheGetProjectById(c, projectId)

	return c.Next()
}

func AddAndEditProjectValidator(c *fiber.Ctx) error {
	pr := new(ProjectRequest)

	if err := c.BodyParser(pr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(pr)
	if err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))

	}

	_, sids, err := subcategory.FindByIds(pr.Subcategory)
	if err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))

	}

	categories, _, err := category.FindByIds(pr.Category)
	if err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))

	}

	finalCategories, err := category.FilterCatesWithSids(categories, sids)
	if err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))

	}

	project := new(Project)
	err = utils.CopyStruct(pr, project)
	if err != nil {
		return errors.Throw(c, err)
	}
	project.Category = finalCategories

	c.Locals("projectBody", project)

	return c.Next()
}
