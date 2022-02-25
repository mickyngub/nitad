package main

import (
	"os"
	"time"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/cronjob"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go.uber.org/zap"
)

//FIXME:  cache
// TODO: cache fiber storage แยก branch

var PORT = os.Getenv("PORT")

func main() {
	utils.InitZap()

	config.Loadenv()
	envErr := config.Checkenv()
	if envErr != nil {
		zap.S().Warn(envErr.Error())
		os.Exit(1)
	}

	database.ConnectDb()
	defer database.DisconnectDb()

	gcp.Init()
	redis.Init()

	app := config.InitApp()

	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Bangkok",
	}))

	app.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
		New:          category.ListCategoryCache(),
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			if path == "/api/v1/category" {
				category.ListCategoryCache(c)
				return true
			}

			return true
		},
	}))

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
		zap.S().Warn("Listen to " + PORT + " Failed!")
		zap.S().Warn("Error: ", err.Error())
	}

}
