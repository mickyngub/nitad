package main

import (
	"log"
	"os"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
)

var PORT = os.Getenv("PORT")

func main() {
	config.Loadenv()
	envErr := config.Checkenv()
	if envErr != nil {
		log.Println(envErr.Error())
		os.Exit(1)
	}

	database.ConnectDb()
	gcp.Init()
	redis.Init()
	app := config.InitApp()
	api.CreateAPI(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Hello, this is NITAD Backend Server v1.3 !"})
	})

	app.All("*", func(c *fiber.Ctx) error {
		return errors.Throw(c, errors.NewNotFoundError("Page"))
	})

	if PORT == "" {
		PORT = "3000"
	}

	log.Println("===== Running on", os.Getenv("APP_ENV"), "stage =====")
	log.Println("===== Listening to port", PORT, "======")

	defer database.DisconnectDb()
	err := app.Listen(":" + PORT)
	if err != nil {
		log.Println("Listen to " + PORT + " Failed!")
		log.Println("Error: ", err.Error())
	}

}
