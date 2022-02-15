package project

import (
	"log"

	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetLookupStage() mongo.Pipeline {
	pipe := mongo.Pipeline{}
	pipe = database.AppendLookupStage(pipe, "category")
	pipe = database.AppendLookupStage(pipe, "subcategory")
	pipe = database.AppendUnsetStage(pipe, "category.subcategory")
	return pipe
}

func IncrementView(id primitive.ObjectID) {
	projectCollection, ctx := database.GetCollection(collectionName)

	_, err := projectCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}},
		},
	)

	// NOTE: logging ??
	if err != nil {
		log.Fatal(err)
	}
}
