package project

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func IncrementView(id primitive.ObjectID, val int) {
	projectCollection, ctx := database.GetCollection(collectionName)

	_, err := projectCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "views", Value: val}}},
		},
	)

	if err != nil {
		zap.S().Warn("Incrementing view error: ", err.Error())
	}
}

func IncrementViewCache(id string, views int) {
	key := "views" + id
	countInt := redis.GetCacheInt(key)
	if countInt != 0 {
		redis.SetCacheInt(key, countInt+1)
		zap.S().Info("incrementing view of ", key, " = ", countInt+1)
		return
	}
	redis.SetCacheInt(key, 1)
}
