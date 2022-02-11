package subcategory

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func isValidObjectId(id string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return objectId, errors.NewBadRequestError("Invalid objectId")
	}
	return objectId, nil
}

func findById(id primitive.ObjectID) (bson.M, errors.CustomError) {
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

func findAll() ([]bson.M, errors.CustomError) {
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
