package spatial

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// get only one spatial in the db
func GetOneSpatial() (Spatial, errors.CustomError) {

	result := Spatial{}
	spatials, err := FindAll()

	if err != nil {
		return result, err
	}

	if len(spatials) == 0 {
		return result, errors.NewBadRequestError("no spatial link yet, please create one first")
	}

	result.ID = spatials[0].ID
	result.Link = spatials[0].Link

	return result, nil
}

func FindAll() ([]Spatial, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	var result []Spatial
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	return result, nil
}

func Add(s *Spatial) (*Spatial, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	insertRes, insertErr := collection.InsertOne(ctx, &s)
	if insertErr != nil {
		return s, errors.NewBadRequestError(insertErr.Error())
	}
	s.ID = insertRes.InsertedID.(primitive.ObjectID)

	return s, nil

}

func Edit(s *Spatial) (*Spatial, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)
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
		return errors.NewBadRequestError("Delete spatial failed!" + err.Error())
	}

	return nil
}
