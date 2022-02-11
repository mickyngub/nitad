package category

import (
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateSubcategoyIds(c *CategoryRequest) ([]primitive.ObjectID, errors.CustomError) {
	subcategoryIds := make([]primitive.ObjectID, len(c.Subcategory))

	for i, subcategoryId := range c.Subcategory {
		objectId, err := functions.IsValidObjectId(subcategoryId)
		if err != nil {
			return subcategoryIds, err

		}

		if _, err = subcategory.FindById(objectId); err != nil {
			return subcategoryIds, err
		}

		subcategoryIds[i] = objectId
	}

	return subcategoryIds, nil
}

func Add(c *CategoryRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}
	subcategoryIds, err := ValidateSubcategoyIds(c)
	if err != nil {
		return result, err
	}

	collection, ctx := database.GetCollection(collectionName)

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: c.Title},
		{Key: "subcategory", Value: subcategoryIds},
	})

	if insertErr != nil {
		return result, errors.NewBadRequestError(insertErr.Error())
	}
	result = map[string]interface{}{
		"id":          insertRes.InsertedID,
		"title":       c.Title,
		"subcategory": subcategoryIds,
	}

	return result, nil
}
