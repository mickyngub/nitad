package subcategory

import (
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindById(oid primitive.ObjectID) (bson.M, errors.CustomError) {
	return database.FindById(oid, collectionName)

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
		{Key: "createdAt", Value: time.Now()},
		{Key: "updatedAt", Value: time.Now()},
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

func Edit(oid primitive.ObjectID, s *Subcategory) (map[string]interface{}, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result map[string]interface{}

	_, updateErr := collection.UpdateByID(
		ctx,
		oid,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: s.Title},
				{Key: "image", Value: s.Image},
				{Key: "updatedAt", Value: time.Now()},
			},
		},
		})

	if updateErr != nil {
		return result, errors.NewBadRequestError(updateErr.Error())
	}

	result = map[string]interface{}{
		"id":    oid,
		"title": s.Title,
		"image": s.Image,
	}

	return result, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete failed!" + err.Error())
	}

	return nil
}
