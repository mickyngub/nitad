package admin

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func IsAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: AuthSuccess,
		ErrorHandler:   AuthError,
		SigningKey:     []byte(JWT_SECRET),
		SigningMethod:  "HS256",
		TokenLookup:    "header:Authorization",
		AuthScheme:     "Bearer",
	})
}

func AuthSuccess(c *fiber.Ctx) error {
	c.Next()
	return nil
}

func AuthError(c *fiber.Ctx, e error) error {
	return errors.Throw(c, errors.NewUnAuthorizedError(e.Error()))
}
