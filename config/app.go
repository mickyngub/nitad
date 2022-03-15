package config

import (
	"os"
	"time"

	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

var app *fiber.App

func GetApp() *fiber.App {
	return app
}

func InitApp() *fiber.App {
	redisStore := redis.Init()

	app = fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Nitad",
		BodyLimit:     100 * 1024 * 1024, // 100 mb
	})

	allowOrigins := "https://nitad-admin-web-app-palmmy2542.vercel.app, https://nitad-web-app.vercel.app"
	if os.Getenv("ALLOW_ORIGINS_ENDPOINT") != "" {
		allowOrigins = os.Getenv("ALLOW_ORIGINS_ENDPOINT")
	}
	zap.S().Info("allow origin ", allowOrigins)
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return errors.Throw(c, errors.NewTooManyRequestsError())
		},
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Bangkok",
	}))

	app.Use(cache.New(cache.Config{
		Expiration: redis.DefaultCacheExpireTime,
		Storage:    redisStore,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Path() + "?" + string(c.Request().URI().QueryString())
		},
		Next: func(c *fiber.Ctx) bool {
			isTrue := project.IsGetProjectPath(c) // handle incrementing view in cache
			return isTrue
		},
	}))

	return app
}
