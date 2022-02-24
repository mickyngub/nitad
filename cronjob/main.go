package cronjob

import (
	"log"

	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Init() {
	c := cron.New()

	c.AddFunc("@every 12h", UpdateProjectViews) // every 12 hours

	c.Start()
}

func UpdateProjectViews() {
	for {
		keys, cursor := redis.FindAllCacheByPrefix("views")

		for _, key := range keys {
			projectId := key[5:]
			objectId, _ := primitive.ObjectIDFromHex(projectId)

			countInt := redis.GetCacheInt(key)

			project.IncrementView(objectId, countInt)

			redis.DeleteCache(key)
			redis.DeleteCache(projectId)
		}

		if cursor == 0 { // no more keys
			break
		}
	}

	log.Println("updating project views")
}
