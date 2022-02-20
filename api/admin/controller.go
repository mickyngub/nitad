package admin

import (
	"strings"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func NewController(
	adminRoute fiber.Router,
) {

	controller := &Controller{}

	//TODO Add auth

	adminRoute.Post("/signup", SignupValidator, controller.Signup)
	adminRoute.Post("/login", LoginValidator, controller.Login)

	adminRoute.Use(IsAuth())
	adminRoute.Get("/profile", controller.Profile)

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

	token, err := CreateToken(&newAdmin)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": fiber.Map{
		"username":     a.Username,
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
	}})

}

// login
func (contc *Controller) Login(c *fiber.Ctx) error {
	a := new(Admin)
	c.BodyParser(a)

	admin, err := FindByUsername(a.Username)
	if err != nil {
		return errors.Throw(c, err)
	}

	if admin == nil {
		return errors.Throw(c, errors.NewInvalidCredentials)
	}

	err = ComparePassword(admin.Password, a.Password)
	if err != nil {
		return errors.Throw(c, err)
	}

	token, err := CreateToken(admin)
	if err != nil {
		return errors.Throw(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": fiber.Map{
		"username":     a.Username,
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
	}})
}

// logout
func (contc *Controller) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": "result"})
}

// get admin profile
func (contc *Controller) Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": fiber.Map{
		"username": claims["username"],
		"exp":      claims["exp"],
	}})

}
