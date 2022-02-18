package admin

import (
	"strings"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	adminRoute fiber.Router,
) {

	controller := &Controller{}

	//TODO Add auth

	adminRoute.Post("/signup", SignupValidator, controller.Signup)
	// adminRoute.Post("/login", LoginValidator, controller.Login)
	//adminRoute.Post("/logout", controller.Logout)

}

type Controller struct{}

// sign up admin
func (contc *Controller) Signup(c *fiber.Ctx) error {
	a := new(AdminSignup)
	c.BodyParser(a)

	a.Username = strings.TrimSpace(a.Username)
	a.Password = strings.TrimSpace(a.Password)
	a.ConfirmPassword = strings.TrimSpace(a.ConfirmPassword)

	if a.Password != a.ConfirmPassword {
		return errors.Throw(c, errors.NewBadRequestError("Password does not match with the confirm password"))
	}

	// hash
	hashedPassword, err := HashPassword(a.Password)
	if err != nil {
		return errors.Throw(c, err)
	}

	admin := Admin{
		Username: a.Username,
		Password: hashedPassword,
	}

	newAdmin, err := CreateAdmin(admin)
	if err != nil {
		return errors.Throw(c, err)
	}

	result, err := CreateJWTToken(newAdmin)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": result})
}

// login
func (contc *Controller) Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "result"})
}

// logout
func (contc *Controller) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "result"})
}
