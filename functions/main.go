package functions

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsValidObjectId(id string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return objectId, errors.NewBadRequestError("Invalid objectId")
	}
	return objectId, nil
}

func RemoveDuplicateObjectIds(oids []primitive.ObjectID) []primitive.ObjectID {
	keys := make(map[primitive.ObjectID]bool)
	list := []primitive.ObjectID{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range oids {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//  Extract files from request body, if no file passed, no error
func ExtractUpdatedFiles(c *fiber.Ctx, key string) ([]*multipart.FileHeader, errors.CustomError) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid input")
	}

	files := form.File[key]
	if len(files) <= 0 {
		return nil, nil
	}

	return files, nil
}

// extractFiles extract files from request body
func ExtractFiles(c *fiber.Ctx, key string) ([]*multipart.FileHeader, errors.CustomError) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid input: " + err.Error())
	}

	files := form.File[key]
	if len(files) <= 0 {
		return nil, errors.NewBadRequestError("at least one file must me provided")
	}

	return files, nil
}

func GetUniqueFilename(filename string) string {
	return fmt.Sprintf("%s-%s.png", time.Now().Format("02-Jan-2006-15:04:05"), strings.TrimSuffix(filename, filepath.Ext(filename)))
}

// for testing purpose only
func WriteFileToPath(f *multipart.FileHeader, filename string) {
	fileContent, _ := f.Open()
	var newErr error
	byteContainer, newErr := ioutil.ReadAll(fileContent)
	filename = fmt.Sprintf("%s.png", strings.TrimSuffix(filename, filepath.Ext(filename)))

	ioutil.WriteFile(filename, byteContainer, 0666)

	if newErr != nil {
		log.Fatal(newErr)
	}
}
