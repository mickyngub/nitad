package category

import (
	"log"

	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewController(
	categoryRoute fiber.Router,
) {

	controller := &Controller{}

	categoryRoute.Get("/", controller.ListCategory)
	categoryRoute.Get("/:categoryId", controller.GetCategory)

	categoryRoute.Use(admin.IsAuth())
	categoryRoute.Post("/", AddAndEditCategoryValidator, controller.AddCategory)
	categoryRoute.Put("/:categoryId", AddAndEditCategoryValidator, controller.EditCategory)
	categoryRoute.Delete("/:categoryId", controller.DeleteCategory)
}

type Controller struct{}

var collectionName = database.COLLECTIONS["CATEGORY"]

// list all categories
func (contc *Controller) ListCategory(c *fiber.Ctx) error {
	var cate []*Category
	key := "allcate"
	err := redis.GetCache(key, &cate)
	if err != nil && err.Error() != "Key does not exist" {
		return errors.Throw(c, err)
	}
	if cate != nil {
		log.Println("getting from cache", key)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": cate})
	}
	categories, err := FindAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	err = redis.SetCache(key, categories)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": categories})
}

// get category by id
func (contc *Controller) GetCategory(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")

	var p []*Category
	err := redis.GetCache(categoryId, &p)
	if err != nil && err.Error() != "Key does not exist" {
		return errors.Throw(c, err)
	}
	if p != nil {
		log.Println("getting from cache", categoryId)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": p})
	}

	objectId, err := functions.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var result Category
	if result, err = GetById(objectId); err != nil {
		return errors.Throw(c, err)
	}
	err = redis.SetCache(categoryId, result)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a category
func (contc *Controller) AddCategory(c *fiber.Ctx) error {
	categoryBody := c.Locals("categoryBody").(*Category)
	sids := c.Locals("sids").([]primitive.ObjectID)

	result, err := Add(categoryBody, sids)

	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

// edit the category
func (contc *Controller) EditCategory(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")
	categoryObjectId, err := functions.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	categoryBody := c.Locals("categoryBody").(*Category)
	sids := c.Locals("sids").([]primitive.ObjectID)

	result, err := Edit(categoryObjectId, categoryBody, sids)

	if err != nil {
		return errors.Throw(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})

}

// delete the category
func (cont *Controller) DeleteCategory(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")
	objectId, err := functions.IsValidObjectId(categoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	err = Delete(objectId)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Delete category " + categoryId + " successfully!"})
}
