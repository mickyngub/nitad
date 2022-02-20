package admin

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindByUsername(username string) (*Admin, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var admin Admin
	err := collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&admin)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError(collectionName)
		} else {
			return nil, errors.NewInternalServerError("getting user error, " + err.Error())
		}
	}

	return &admin, nil
}

func CreateAdmin(a Admin) (Admin, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	admin, err := FindByUsername(a.Username)
	if err != nil && err.Error() != errors.NewNotFoundError(collectionName).Error() {
		return a, err
	}

	if admin != nil {
		return a, errors.NewBadRequestError("Username has already been taken")
	}

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "username", Value: a.Username},
		{Key: "password", Value: a.Password},
	})

	if insertErr != nil {
		return a, errors.NewBadRequestError(insertErr.Error())
	}

	a.ID = insertRes.InsertedID.(primitive.ObjectID)
	return a, nil
}
