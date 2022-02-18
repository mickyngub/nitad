package admin

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAdmin(a Admin) (Admin, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	// TODO: check for existing one
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
