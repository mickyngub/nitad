package category

import (
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindById(id primitive.ObjectID) (bson.M, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	matchStage := bson.D{{"$match", bson.D{{"_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "subcategory"}, {"localField", "subcategory"}, {"foreignField", "_id"}, {"as", "subcategory"}}}}

	cursor, err := categoryCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage})
	var result []bson.M
	if err != nil {
		return result[0], errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &result); err != nil {
		return result[0], errors.NewBadRequestError(err.Error())
	}

	return result[0], nil
}

func FindAll() ([]bson.M, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "subcategory"}, {"localField", "subcategory"}, {"foreignField", "_id"}, {"as", "subcategory"}}}}

	cursor, err := categoryCollection.Aggregate(ctx, mongo.Pipeline{lookupStage})
	var result []bson.M
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())

	}
	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	return result, nil
}

func Add(c *CategoryRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}
	subcategoryIds, err := subcategory.ValidateIds(c.Subcategory)
	if err != nil {
		return result, err
	}

	collection, ctx := database.GetCollection(collectionName)

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: c.Title},
		{Key: "subcategory", Value: subcategoryIds},
	})

	if insertErr != nil {
		return result, errors.NewBadRequestError(insertErr.Error())
	}
	result = map[string]interface{}{
		"id":          insertRes.InsertedID,
		"title":       c.Title,
		"subcategory": subcategoryIds,
	}

	return result, nil
}
