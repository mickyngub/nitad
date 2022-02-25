package main

import (
	"os"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/cronjob"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/logger"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

//FIXME: log / cache
// TODO: cache fiber storage แยก branch

var PORT = os.Getenv("PORT")

func main() {
	logger.InitZap()

	config.Loadenv()
	envErr := config.Checkenv()
	if envErr != nil {
		zap.S().Fatal(envErr.Error())
		os.Exit(1)
	}

	database.ConnectDb()
	defer database.DisconnectDb()

	gcp.Init()
	redis.Init()
	app := config.InitApp()
	api.CreateAPI(app)
	cronjob.Init()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Hello, this is NITAD Backend Server v1.7 !"})
	})

	app.All("*", func(c *fiber.Ctx) error {
		return errors.Throw(c, errors.NewNotFoundError("Page"))
	})

	if PORT == "" {
		PORT = "3000"
	}

	zap.S().Info("===== Running on ", os.Getenv("APP_ENV"), " stage =====")
	zap.S().Info("===== Listening to port ", PORT, " ======")

	err := app.Listen(":" + PORT)
	if err != nil {
		zap.S().Fatal("Listen to " + PORT + " Failed!")
		zap.S().Fatal("Error: ", err.Error())
	}

}
