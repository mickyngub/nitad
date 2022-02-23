package api

import (
	"github.com/birdglove2/nitad-backend/api/admin"
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/api/spatial"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/gofiber/fiber/v2"
)

const API_PREFIX = "/api/v1"

func CreateAPI(app *fiber.App) {
	v1 := app.Group(API_PREFIX)

	subcategory.NewController(v1.Group("/subcategory"))
	category.NewController(v1.Group("/category"))
	project.NewController(v1.Group("/project"))

	spatial.NewController(v1.Group("/spatial"))

	admin.NewController(v1.Group("/admin"))

}
