package config

import "github.com/gofiber/fiber/v2"

var app *fiber.App

func GetApp() *fiber.App {
	return app
}

func InitApp() *fiber.App {
	app = fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Nitad",
	})
	return app
}
