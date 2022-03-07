package setup

import (
	context "context"
	"image"
	multipart "mime/multipart"
	"os"
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var subcateRepo subcategory.Repository
var cateRepo category.Repository
var projectRepo project.Repository
var app *fiber.App

func TestMain(m *testing.M) {
	config.Loadenv()

	client := database.ConnectDb(os.Getenv("MONGO_URI"))

	subcateRepo = subcategory.NewRepository(client)
	cateRepo = category.NewRepository(client)
	projectRepo = project.NewRepository(client)

	os.Exit(m.Run())
}

func NewTestApp(t *testing.T) *fiber.App {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gcpService := NewMockUploader(ctrl)

	app = fiber.New()
	api.CreateAPI(app, gcpService)
	return app
}

func AddMockSubcategory(t *testing.T) *subcategory.Subcategory {
	dummySubcate := subcategory.Subcategory{
		Title: "dummy subcate title",
		Image: "dummy subcate image url",
	}

	adddedSubcategory, err := subcateRepo.AddSubcategory(context.Background(), &dummySubcate)
	require.Equal(t, err, nil)
	require.Equal(t, dummySubcate.Title, adddedSubcategory.Title)
	require.Equal(t, dummySubcate.Image, adddedSubcategory.Image)
	require.NotEqual(t, nil, adddedSubcategory.ID)

	return adddedSubcategory
}

func AddMockCategory(t *testing.T, subcate *subcategory.Subcategory) *category.Category {
	dummyCate := category.Category{
		Title: "dummy cate title",
	}

	addedCategory, err := category.Add(&dummyCate, []primitive.ObjectID{subcate.ID})
	require.Equal(t, err, nil)
	require.Equal(t, dummyCate.Title, addedCategory.Title)
	require.Equal(t, dummyCate.Subcategory, addedCategory.Subcategory)
	require.NotEqual(t, nil, addedCategory.ID)

	return addedCategory
}

func AddMockProject(t *testing.T, cate *category.Category) *project.Project {
	dummyProj := project.Project{
		Title:       "dummy proj title",
		Description: "dumym proj description",
		Authors:     []string{"dumym proj Authors"},
		Emails:      []string{"dumym proj Emails"},
		Inspiration: "dumym proj Inspiration",
		Abstract:    "dumym proj Abstract",
		Images:      []string{"dumym proj Images"},
		Videos:      []string{"dumym proj Videos"},
		Keywords:    []string{"dumym proj Keywords"},
		Report:      "dumym proj Report",
		VirtualLink: "dumym proj VirtualLink",
		Status:      "dumym proj Status",
		Category:    []category.Category{*cate},
	}
	addedProject, err := project.Add(&dummyProj)
	require.Equal(t, err, nil)
	require.Equal(t, dummyProj.Title, addedProject.Title)
	require.Equal(t, dummyProj.Category, addedProject.Category)
	require.NotEqual(t, nil, addedProject.ID)
	return addedProject
}

func DeleteMock(t *testing.T, proj *project.Project, cate *category.Category, subcate *subcategory.Subcategory) {
	err := subcateRepo.DeleteSubcategory(context.Background(), subcate.ID)
	require.Nil(t, err, "Delete subcate failed")

	err = cateRepo.DeleteCategory(context.Background(), subcate.ID)
	require.Nil(t, err, "Delete cate failed")

	err = projectRepo.DeleteProject(context.Background(), subcate.ID)
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

func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
