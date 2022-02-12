package project

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
)

func Add(c *ProjectRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}
	subcategoryIds, err := subcategory.ValidateIds(c.Subcategory)
	if err != nil {
		return result, err
	}

	categoryIds, err := category.ValidateIds(c.Category)
	if err != nil {
		return result, err
	}

	subcategoryIds = functions.RemoveDuplicateObjectIds(subcategoryIds)
	categoryIds = functions.RemoveDuplicateObjectIds(categoryIds)

	collection, ctx := database.GetCollection(collectionName)

	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: c.Title},
		{Key: "description", Value: c.Description},
		{Key: "authors", Value: c.Authors},
		{Key: "emails", Value: c.Emails},
		{Key: "inspiration", Value: c.Inspiration},
		{Key: "abstract", Value: c.Abstract},
		{Key: "images", Value: c.Images},
		{Key: "videos", Value: c.Videos},
		{Key: "keywords", Value: c.Keywords},
		{Key: "category", Value: categoryIds},
		{Key: "subcategory", Value: subcategoryIds},
	})

	if insertErr != nil {
		return result, errors.NewBadRequestError(insertErr.Error())
	}

	result = map[string]interface{}{
		"id":          insertRes.InsertedID,
		"title":       c.Title,
		"description": c.Description,
		"authors":     c.Authors,
		"emails":      c.Emails,
		"inspiration": c.Inspiration,
		"abstract":    c.Abstract,
		"images":      c.Images,
		"videos":      c.Videos,
		"keywords":    c.Keywords,
		"category":    categoryIds,
		"subcategory": subcategoryIds,
	}

	return result, nil
}
