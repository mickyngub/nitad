package subcategory

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindById(id primitive.ObjectID) (bson.M, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result bson.M
	err := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return result, errors.NewNotFoundError("subcategoryId")
		} else {
			return result, errors.NewBadRequestError("something went wrong")
		}
	}

	return result, nil
}

func FindAll() ([]bson.M, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result []bson.M
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	return result, nil
}

func Add(s *Subcategory) (map[string]interface{}, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result map[string]interface{}

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: s.Title},
		{Key: "image", Value: s.Image},
	})

	if insertErr != nil {
		return result, errors.NewBadRequestError(insertErr.Error())
	}

	result = map[string]interface{}{
		"id":    insertRes.InsertedID,
		"title": s.Title,
		"image": s.Image,
	}

	return result, nil
}
