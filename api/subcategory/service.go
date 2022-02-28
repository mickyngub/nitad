package subcategory

import (
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// get the SUBCATEGORY from requested id
// ** different from FindById
func GetById(oid primitive.ObjectID) (Subcategory, errors.CustomError) {
	m, err := database.GetElementById(oid, collectionName)
	return BsonToSubcategory(m), err
}

func FindAll() ([]Subcategory, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result []Subcategory
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	return result, nil
}

func Add(s *Subcategory) (*Subcategory, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	insertRes, insertErr := collection.InsertOne(ctx, &s)
	if insertErr != nil {
		return s, errors.NewBadRequestError(insertErr.Error())
	}

	s.ID = insertRes.InsertedID.(primitive.ObjectID)

	return s, nil
}

func Edit(s *Subcategory) (*Subcategory, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	now := time.Now()
	s.UpdatedAt = now
	_, updateErr := collection.UpdateByID(
		ctx,
		s.ID,
		bson.D{{
			Key: "$set", Value: &s}})

	if updateErr != nil {
		return s, errors.NewBadRequestError(updateErr.Error())
	}

	return s, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete subcategory failed!" + err.Error())
	}

	return nil
}
