package main

import (
	"os"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/cronjob"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	utils.InitZap()

	config.Loadenv()
	zap.S().Info("===== Running on ", os.Getenv("APP_ENV"), " stage =====")
	database.ConnectDb(os.Getenv("MONGO_URI"))
	defer database.DisconnectDb()

	uploader := gcp.Init()

	app := config.InitApp()

	api.CreateAPI(app, uploader)

	cronjob.Init()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Hello, this is NITAD Backend " + os.Getenv("APP_ENV") + " Server v2.8  !"})
	})

	app.All("*", func(c *fiber.Ctx) error {
		return errors.Throw(c, errors.NewNotFoundError("Page"))
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	zap.S().Info("===== Listening to port ", PORT, " ======")

	err := app.Listen(":" + PORT)
	if err != nil {
		zap.S().Warn("Listen to " + PORT + " Failed!")
		zap.S().Warn("Error: ", err.Error())
	}
}
