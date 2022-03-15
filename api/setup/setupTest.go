package setup

import (
	"bytes"
	context "context"
	"io"
	multipart "mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var SubcateRepo subcategory.Repository
var CateRepo category.Repository
var ProjectRepo project.Repository
var app *fiber.App

var Token string

func NewTestApp(t *testing.T) (*fiber.App, *MockUploader) {
	config.Loadenv()

	client := database.ConnectDb(os.Getenv("MONGO_URI"))
	SubcateRepo = subcategory.NewRepository(client)
	CateRepo = category.NewRepository(client)
	ProjectRepo = project.NewRepository(client)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gcpService := NewMockUploader(ctrl)

	app = fiber.New()

	api.CreateAPI(app, gcpService)

	Token = Login(t, app)

	return app, gcpService
}

func AddMockSubcategory(t *testing.T) *subcategory.Subcategory {
	dummySubcate := subcategory.Subcategory{
		Title: "dummy subcate title",
		Image: "dummy subcate image url",
	}

	adddedSubcategory, err := SubcateRepo.AddSubcategory(context.Background(), &dummySubcate)
	require.Nil(t, err)
	require.Equal(t, dummySubcate.Title, adddedSubcategory.Title)
	require.Equal(t, dummySubcate.Image, adddedSubcategory.Image)
	require.NotNil(t, adddedSubcategory.ID)

	return adddedSubcategory
}

func AddMockCategory(t *testing.T, subcate *subcategory.Subcategory) *category.Category {
	dummyCate := category.CategoryDTO{
		Title: "dummy cate title",
	}

	dummyCate.Subcategory = []string{subcate.ID.Hex()}
	addedCategory, err := CateRepo.AddCategory(context.Background(), &dummyCate)
	require.Equal(t, err, nil)
	require.Equal(t, dummyCate.Title, addedCategory.Title)
	require.Equal(t, dummyCate.Subcategory, addedCategory.Subcategory)
	require.NotEqual(t, nil, addedCategory.ID)

	cate, err := CateRepo.GetCategoryById(context.Background(), addedCategory.ID)
	require.Equal(t, err, nil)

	return cate
}

func AddMockProject(t *testing.T, cate *category.Category) *project.Project {
	dummyProj := project.Project{
		Title:       "dummy proj title",
		Description: "dummy proj description",
		Authors:     []string{"dummy proj Authors"},
		Emails:      []string{"dummy proj Emails"},
		Inspiration: "dummy proj Inspiration",
		Abstract:    "dummy proj Abstract",
		Images:      []string{"project/images/dummy proj images"},
		Videos:      []string{"dummy proj Videos"},
		Keywords:    []string{"dummy proj Keywords"},
		Report:      "project/report/dummy proj report",
		VirtualLink: "dummy proj VirtualLink",
		Status:      "dummy proj Status",
		Category:    []category.Category{*cate},
	}

	addedProject, err := ProjectRepo.AddProject(context.Background(), &dummyProj)
	require.Equal(t, err, nil)
	require.Equal(t, dummyProj.Title, addedProject.Title)
	require.Equal(t, dummyProj.Category, addedProject.Category)
	require.NotEqual(t, nil, addedProject.ID)
	return addedProject
}

func DeleteMock(t *testing.T, proj *project.Project, cate *category.Category, subcate *subcategory.Subcategory) {
	DeleteMockSubcategory(t, subcate)
	DeleteMockCategory(t, cate)
	DeleteMockProject(t, proj)
}

func DeleteMockSubcategory(t *testing.T, subcate *subcategory.Subcategory) {
	err := SubcateRepo.DeleteSubcategory(context.Background(), subcate.ID)
	require.Nil(t, err, "Delete subcate failed")
}

func DeleteMockCategory(t *testing.T, cate *category.Category) {
	err := CateRepo.DeleteCategory(context.Background(), cate.ID)
	require.Nil(t, err, "Delete cate failed")
}

func DeleteMockProject(t *testing.T, proj *project.Project) {
	err := ProjectRepo.DeleteProject(context.Background(), proj.ID)
	require.Nil(t, err, "Delete proj failed")
}

func RandomImages(n int) []*multipart.FileHeader {
	results := make([]*multipart.FileHeader, n)
	for i := 0; i < n; i++ {
		results[i] = &multipart.FileHeader{
			Filename: utils.RandomString(5) + ".jpg",
		}
	}
	return results
}

func OpenFileFromPath(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func Upload(method string, url string, values map[string]interface{}) (req *http.Request, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, val := range values {
		// if file upload as file
		if key == "image" || key == "images" || key == "report" {
			if osFile, ok := val.(*os.File); ok {
				fw, err := w.CreateFormFile(key, osFile.Name())
				if err != nil {
					return nil, err
				}

				if _, err = io.Copy(fw, val.(io.Reader)); err != nil {
					return nil, err
				}
			}

		} else {

			// if slice, use for loop to add
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				for _, v := range val.([]string) {
					fw, err := w.CreateFormField(key)
					if err != nil {
						return nil, err
					}
					if _, err = io.Copy(fw, strings.NewReader(v)); err != nil {
						return nil, err
					}
				}

				// others are all string type
			} else {
				fw, err := w.CreateFormField(key)
				if err != nil {
					return nil, err
				}
				if _, err = io.Copy(fw, strings.NewReader(val.(string))); err != nil {
					return nil, err
				}
			}
		}
	}
	w.Close()

	req, err = http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil

}

func CreateMultipartFormDataRequest(method string, url string, values map[string]io.Reader) (req *http.Request, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return nil, err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		}

		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err = http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}
