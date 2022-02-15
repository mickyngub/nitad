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

func FindById(id primitive.ObjectID) (bson.M, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	matchStage := bson.D{{"$match", bson.D{{"_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "subcategory"}, {"localField", "subcategory"}, {"foreignField", "_id"}, {"as", "subcategory"}}}}

	cursor, err := categoryCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage})
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
		"id":          insertRes.InsertedID,
		"title":       c.Title,
		"subcategory": subcategoryIds,
	}

	return result, nil
}

func Edit(oid primitive.ObjectID, c *CategoryRequest) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)
	// oldCategory, err := database.FindById(oid, collectionName)
	// if err != nil {
	// 	return err
	// }

	// oc := BsonToCategory(oldCategory)
	// sidsString := append(c.Subcategory, oc.Subcategory...)

	subcategoryIds, err := subcategory.ValidateIds(c.Subcategory)
	if err != nil {
		return err
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
		return errors.NewBadRequestError(updateErr.Error())
	}
	return nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete failed!" + err.Error())
	}

	return nil
}
