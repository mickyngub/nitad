package subcategory

import (
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindById(id primitive.ObjectID) (bson.M, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result bson.M
	err := collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return result, errors.NewNotFoundError("subcategoryId")
		} else {
			return result, errors.NewBadRequestError(err.Error())
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

// validate requested string of subcategoryIds
// and return valid []objectId, otherwise error
func ValidateIds(sids []string) ([]primitive.ObjectID, errors.CustomError) {
	objectIds := make([]primitive.ObjectID, len(sids))

	for i, sid := range sids {
		objectId, err := ValidateId(sid)
		if err != nil {
			return objectIds, err
		}

		objectIds[i] = objectId
	}

	return objectIds, nil
}

// validate requested string of a single subcategoryId
// and return valid objectId, otherwise error
func ValidateId(sid string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := functions.IsValidObjectId(sid)
	if err != nil {
		return objectId, err

	}
	// if err != nil >> id is not found
	if _, err = FindById(objectId); err != nil {
		return objectId, err
	}

	return objectId, nil
}
