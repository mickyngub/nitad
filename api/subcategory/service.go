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
	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: s.Title},
		{Key: "image", Value: s.Image},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	})

	if insertErr != nil {
		return s, errors.NewBadRequestError(insertErr.Error())
	}

	s.ID = insertRes.InsertedID.(primitive.ObjectID)
	s.CreatedAt = now
	s.UpdatedAt = now

	return s, nil
}

func Edit(oid primitive.ObjectID, ns *Subcategory) (*Subcategory, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	result := new(Subcategory)

	now := time.Now()
	_, updateErr := collection.UpdateByID(
		ctx,
		oid,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: ns.Title},
				{Key: "image", Value: ns.Image},
				{Key: "updatedAt", Value: now},
			},
		},
		})

	if updateErr != nil {
		return result, errors.NewBadRequestError(updateErr.Error())
	}

	result.ID = oid
	result.Title = ns.Title
	result.Image = ns.Image
	result.CreatedAt = ns.CreatedAt
	result.UpdatedAt = now

	return result, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete subcategory failed!" + err.Error())
	}

	return nil
}
