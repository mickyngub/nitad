package spatial

import (
	"strconv"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	spatialRoute fiber.Router,
) {

	controller := &Controller{}

	spatialRoute.Get("/:roomNumber", controller.GetLink)
}

type Controller struct{}

var links = []string{
	"https://google.com",
	"https://facebook.com",
	"https://youtube.com",
	"https://google.com",
	"https://facebook.com",
	"https://youtube.com",
	"https://google.com",
	"https://facebook.com",
	"https://google.com",
	"https://facebook.com",
}

// get spatial new link
func (contc *Controller) GetLink(c *fiber.Ctx) error {
	roomNumberStr := c.Params("roomNumber")
	roomNumber, err := strconv.Atoi(roomNumberStr)
	if err != nil || roomNumber < 1 || roomNumber > 10 {
		return errors.Throw(c, errors.NewBadRequestError("room number only valid from 1 to 10"))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result": fiber.Map{
		"roomNumber": roomNumber,
		"link":       links[roomNumber-1],
	}})
}
