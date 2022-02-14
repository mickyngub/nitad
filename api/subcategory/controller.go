package subcategory

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewController(
	subcategoryRoute fiber.Router,
) {

	controller := &Controller{}

	subcategoryRoute.Get("/", controller.ListSubcategory)
	subcategoryRoute.Get("/:subcategoryId", controller.GetSubcategory)

	//TODO add AUTH for POST/PUT/DELETE

	subcategoryRoute.Post("/", controller.AddSubcategory)
	// subcategoryRoute.Put("/:subcategoryId", controller.EditSubcategory)
	// subcategoryRoute.Delete("/:subcategoryId", controller.DeleteSubcategory)
}

type Controller struct {
	// service Service
}

var collectionName = database.COLLECTIONS["SUBCATEGORY"]

// list all subcategories
func (contc *Controller) ListSubcategory(c *fiber.Ctx) error {
	subcategories, err := FindAll()
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": subcategories})
}

// get subcategory by id
func (contc *Controller) GetSubcategory(c *fiber.Ctx) error {
	subcategoryId := c.Params("subcategoryId")

	objectId, err := functions.IsValidObjectId(subcategoryId)
	if err != nil {
		return errors.Throw(c, err)
	}

	var result bson.M
	if result, err = FindById(objectId); err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// add a subcategory
func (contc *Controller) AddSubcategory(c *fiber.Ctx) error {

	p := new(Subcategory)

	//TODO: handle this bodyParser middleware
	if err := c.BodyParser(p); err != nil {
		return errors.Throw(c, errors.NewBadRequestError(err.Error()))
	}

	files, err := contc.extractFiles(c, "image")
	if err != nil {
		return errors.Throw(c, err)
	}

	//TODO: handle file upload
	imageURLs := UploadImage(files)

	// var subcate Subcategory
	p.Image = imageURLs[0]

	result, err := Add(p)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

//TODO: upload to AWS
func UploadImage(files []*multipart.FileHeader) []string {

	for _, file := range files {
		// should get image url instead
		// writing file here is just for testing purpose
		WriteFile(file, file.Filename)
	}

	dummyURL := "https://www.einfochips.com/blog/wp-content/uploads/2018/11/how-to-develop-machine-learning-applications-for-business-featured.jpg"
	return []string{dummyURL}
}

func WriteFile(f *multipart.FileHeader, filename string) {
	fileContent, _ := f.Open()
	var newErr error
	byteContainer, newErr := ioutil.ReadAll(fileContent)
	filename = fmt.Sprintf("%s.png", strings.TrimSuffix(filename, filepath.Ext(filename)))

	ioutil.WriteFile(filename, byteContainer, 0666)

	if newErr != nil {
		log.Fatal(newErr)
	}
}

// // edit the subcategory
// func (contc *Controller) EditSubcategory(c *fiber.Ctx) error {}

// // delete the subcategory
// func (cont *Controller) DeleteSubcategory(c *fiber.Ctx) error {}
