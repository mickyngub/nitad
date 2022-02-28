package utils

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetUniqueFilename(filename string) (string, string) {
	// return fmt.Sprintf("%s-%s.png", time.Now().Format("02-Jan-2006-15:04:05"), strings.TrimSuffix(filename, filepath.Ext(filename)))
	filetype := "images"
	if filepath.Ext(filename) == ".pdf" {
		filetype = "reports"
	}
	return fmt.Sprintf("%s-%s", time.Now().Format("02-Jan-2006-15:04:05"), filename), filetype
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
		return nil, errors.NewBadRequestError("Invalid form input: " + err.Error())
	}

	files := form.File[key]
	if len(files) <= 0 {
		return nil, errors.NewBadRequestError("At least one file must me provided for " + key)
	}

	return files, nil
}

// for testing purpose only
func WriteFileToPath(f *multipart.FileHeader, filename string) {
	fileContent, _ := f.Open()
	var newErr error
	byteContainer, newErr := ioutil.ReadAll(fileContent)
	filename = fmt.Sprintf("%s.png", strings.TrimSuffix(filename, filepath.Ext(filename)))

	ioutil.WriteFile(filename, byteContainer, 0666)

	if newErr != nil {
		zap.S().Warn(newErr.Error())
	}
}
