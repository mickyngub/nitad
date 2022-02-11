package subcategory

import (
	"fmt"
	"log"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewController(
	subcategoryRoute fiber.Router,
) {

	controller := &Controller{}

	subcategoryRoute.Get("/", controller.Listsubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.Getsubcategory)

	//TODO add AUTH for POST/PUT/DELETE

	subcategoryRoute.Post("/", controller.Addsubcategory)
	// subcategoryRoute.Put("/:subcategoryId", controller.editsubcategory)
	// subcategoryRoute.Delete("/:subcategoryId", controller.deletesubcategory)
}

type Controller struct {
	// service Service
}

var collectionName = database.COLLECTIONS["SUBCATEGORY"]

// Get subcategory by id
func (contc *Controller) Getsubcategory(c *fiber.Ctx) error {

	collection, ctx := database.GetCollection(collectionName)

	id, err := primitive.ObjectIDFromHex(c.Params("subcategoryId"))
	if err != nil {
		return errors.Throw(c, errors.NewNotFoundError("id is in valid 2"))
	}

	var result bson.M
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return errors.Throw(c, errors.NewNotFoundError("subcategoryId"))
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": result})

}

// list all subcategories
func (contc *Controller) Listsubcategory(c *fiber.Ctx) error {

	collection, ctx := database.GetCollection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var subcategories []bson.M
	if err = cursor.All(ctx, &subcategories); err != nil {
		log.Fatal(err)
	}

	for i, subcategory := range subcategories {
		fmt.Println(i+1, subcategory)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": subcategories})
}

// add a subcategory
func (contc *Controller) Addsubcategory(c *fiber.Ctx) error {

	p := new(subcategory)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	collection, ctx := database.GetCollection(collectionName)

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: p.Title},
		{Key: "images", Value: p.Images},
	})
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	res := map[string]interface{}{
		"id":     insertRes.InsertedID,
		"title":  p.Title,
		"images": p.Images,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": res})

}

// // edit the subcategory
// func (contc *Controller) editsubcategory(c *fiber.Ctx) error {}

// // delete the subcategory
// func (cont *Controller) deletesubcategory(c *fiber.Ctx) error {}
