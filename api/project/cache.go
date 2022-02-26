package project

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/birdglove2/nitad-backend/redis"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func IsGetProjectPath(c *fiber.Ctx) bool {
	path := c.Path()
	if strings.Contains(path, "/api/v1/project") {
		// log.Println("2", c.Params("projectId")) //WARNING :wtf  ไม่เจอ งง ????
		pathArr := strings.Split(path, "/")
		projectId := pathArr[len(pathArr)-1]
		_, err := utils.IsValidObjectId(projectId)
		if err != nil {
			return false
		}
		if len(projectId) > 0 {
			// log.Println("2")
			return true
		}
	}
	return false
}

func HandleCacheGetProjectById(c *fiber.Ctx, projectId string) *Project {
	var p Project
	b, _ := redis.GetStore().Get(c.Path())
	if len(b) > 0 {
		IncrementViewCache(projectId, p.Views)
		err := json.Unmarshal(b, &p)
		if err != nil {
			zap.S().Warn("Unmarshal fail", err.Error())
		}
		log.Println("not nexting")

		return &p
	}
	log.Println("nexting")
	return nil
}
