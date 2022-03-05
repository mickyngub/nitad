package api

import (
	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/connection"
	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/api/search"
	"github.com/birdglove2/nitad-backend/api/spatial"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
)

const API_PREFIX = "/api/v1"

func CreateAPI(app *fiber.App, gcpService gcp.Uploader) {
	client := database.GetClient()

	subcategoryRepository := subcategory.NewRepository(client)
	categoryRepository := category.NewRepository(client)

	subcategoryService := subcategory.NewService(subcategoryRepository, gcpService)
	categoryService := category.NewService(categoryRepository, subcategoryService)

	v1 := app.Group(API_PREFIX)
	search.NewController(v1.Group("/search"))
	connection.NewController(v1.Group("/connection"))
	subcategory.NewController(subcategoryService, v1.Group("/subcategory"))
	category.NewController(categoryService, v1.Group("/category"))
	project.NewController(gcpService, v1.Group("/project"))
	spatial.NewController(v1.Group("/spatial"))
	admin.NewController(v1.Group("/admin"))
}
