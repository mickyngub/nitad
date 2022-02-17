package category

import (
	"time"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindById(oid primitive.ObjectID) (bson.M, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	pipe := mongo.Pipeline{}
	pipe = database.AppendMatchIdStage(pipe, "_id", oid)
	pipe = database.AppendLookupStage(pipe, "subcategory")

	cursor, err := categoryCollection.Aggregate(ctx, pipe)
	var result []bson.M
	if err != nil {
		return bson.M{}, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &result); err != nil {
		return bson.M{}, errors.NewBadRequestError(err.Error())
	}

	if len(result) == 0 {
		return bson.M{}, errors.NewNotFoundError("categoryId")

	}

	return result[0], nil
}

func FindAll() ([]bson.M, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	pipe := mongo.Pipeline{}
	pipe = database.AppendLookupStage(pipe, "subcategory")

	cursor, err := categoryCollection.Aggregate(ctx, pipe)
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

	subcategoryIds = functions.RemoveDuplicateObjectIds(subcategoryIds)

	collection, ctx := database.GetCollection(collectionName)

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: c.Title},
		{Key: "subcategory", Value: subcategoryIds},
		{Key: "createdAt", Value: time.Now()},
		{Key: "updatedAt", Value: time.Now()},
	})

	if insertErr != nil {
		return result, errors.NewBadRequestError(insertErr.Error())
	}
	result = map[string]interface{}{
		"_id":         insertRes.InsertedID,
		"title":       c.Title,
		"subcategory": subcategoryIds,
	}

	return result, nil
}

func Edit(oid primitive.ObjectID, c *CategoryRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}

	collection, ctx := database.GetCollection(collectionName)

	subcategoryIds, err := subcategory.ValidateIds(c.Subcategory)
	if err != nil {
		return result, err
	}
	subcategoryIds = functions.RemoveDuplicateObjectIds(subcategoryIds)

	_, updateErr := collection.UpdateByID(
		ctx,
		oid,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: c.Title},
				{Key: "subcategory", Value: subcategoryIds},
				{Key: "updatedAt", Value: time.Now()},
			},
		},
		})

	if updateErr != nil {
		return result, errors.NewBadRequestError(updateErr.Error())
	}

	result = map[string]interface{}{
		"_id":         c.Title,
		"title":       c.Title,
		"subcategory": subcategoryIds,
	}
	return result, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete category failed!" + err.Error())
	}

	return nil
}
