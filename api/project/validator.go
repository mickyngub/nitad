package project

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/api/validators"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func AddProjectValidator(c *fiber.Ctx) error {
	pr := new(ProjectRequest)

	if err := c.BodyParser(pr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*pr)
	if err != nil {
		return errors.Throw(c, err)
	}

	subcategories, sids, err := subcategory.FindByIds(pr.Subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	category, err := category.FindById(pr.Category)

	if err != nil {
		return errors.Throw(c, err)
	}

	project := new(Project)
	project.Title = pr.Title
	project.Description = pr.Description
	project.Authors = pr.Authors
	project.Emails = pr.Emails
	project.Inspiration = pr.Inspiration
	project.Abstract = pr.Abstract
	project.Videos = pr.Videos
	project.Keywords = pr.Keywords
	project.Subcategory = subcategories
	project.Category = category

	c.Locals("projectBody", project)
	c.Locals("cid", category.ID)
	c.Locals("sids", sids)

	return c.Next()
}

func EditProjectValidator(c *fiber.Ctx) error {
	upr := new(UpdateProjectRequest)

	if err := c.BodyParser(upr); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	err := validators.ValidateStruct(*upr)
	if err != nil {
		return errors.Throw(c, err)
	}

	subcategories, sids, err := subcategory.FindByIds(upr.Subcategory)
	if err != nil {
		return errors.Throw(c, err)
	}

	category, err := category.FindById(upr.Category)

	if err != nil {
		return errors.Throw(c, err)
	}

	project := new(UpdateProject)
	project.Title = upr.Title
	project.Description = upr.Description
	project.Authors = upr.Authors
	project.Emails = upr.Emails
	project.Inspiration = upr.Inspiration
	project.Abstract = upr.Abstract
	project.Videos = upr.Videos
	project.Keywords = upr.Keywords
	project.Subcategory = subcategories
	project.Category = category
	project.DeleteImages = upr.DeleteImages

	c.Locals("updateProjectBody", project)
	c.Locals("cid", category.ID)
	c.Locals("sids", sids)

	return c.Next()
}
