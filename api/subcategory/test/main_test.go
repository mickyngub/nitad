package subcategory_test

import (
	"fmt"
	"testing"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
)

func newTestServer(t *testing.T, gcpService gcp.ClientUploader) *fiber.App {
	fmt.Println("newTestServer", 0)

	database.ConnectDb(mongo_URI)
	fmt.Println("newTestServer", 1)
	app := fiber.New()
	fmt.Println("newTestServer", 2)
	api.CreateAPI(app, gcpService)

	fmt.Println("newTestServer", 3)
	return app
}
