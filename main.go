package main

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func getApp() *fiber.App {
	return app
}

const API_PREFIX = "/api/v1"

func main() {
	config.Loadenv() // load env

	database.ConnectDb() // connect database mongoDB

	// initApp
	app = fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Nitad",
	})

	v1 := app.Group(API_PREFIX)

	subcategory.NewController(v1.Group("/subcategory"))
	category.NewController(v1.Group("/category"))
	// project.NewController(v1.Group("/project"))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": "Hello World!"})

	})

	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"result": "Page Not Found"})
	})

	app.Listen(":3000")

}
