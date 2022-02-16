package database

import (
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindById(oid primitive.ObjectID, collectionName string) (bson.M, errors.CustomError) {
	collection, ctx := GetCollection(collectionName)

	var result bson.M
	err := collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return result, errors.NewNotFoundError(collectionName)
		} else {
			return result, errors.NewBadRequestError(err.Error())
		}
	}
	return result, nil
}
