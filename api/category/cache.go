package category

import (
	"log"

	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
)

func ListCategoryCache(c *fiber.Ctx) error {
	log.Println("check 2")

	var cate []*Category
	key := "allcate"
	redis.GetCache(key, &cate)
	if cate != nil {
		log.Println("cache cate success")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "result c": cate})
	}
	log.Println("cache cate not found")

	c.Next()
	return nil
}
