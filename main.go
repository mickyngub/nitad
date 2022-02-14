package main

import (
	"log"
	"os"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

var PORT = os.Getenv("PORT")

func main() {

	config.Loadenv()
	database.ConnectDb()
	app := config.InitApp()
	api.CreateAPI(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "Hello, this is NITAD Backend Server!"})
	})

	app.All("*", func(c *fiber.Ctx) error {
		return errors.Throw(c, errors.NewNotFoundError("Page"))
	})

	if PORT == "" {
		PORT = "3000"
	}

	PORT = ":" + PORT

	log.Println("Listening to ", PORT)

	err := app.Listen(PORT)
	if err != nil {
		log.Printf("Listen to %s Failed", PORT)
		log.Fatal("Error: ", err.Error())
	}

}
