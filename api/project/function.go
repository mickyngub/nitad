package project

import (
	"log"
	"time"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetLookupStage() []bson.D {
	lookupCategoryStage := bson.D{{"$lookup", bson.D{{"from", "category"},
		{"localField", "category"}, {"foreignField", "_id"}, {"as", "category"}}}}

	unsetStage := bson.D{{"$unset", "category.subcategory"}}

	lookupSubcategoryStage := bson.D{{"$lookup", bson.D{{"from", "subcategory"}, {"localField", "subcategory"}, {"foreignField", "_id"}, {"as", "subcategory"}}}}

	return []bson.D{lookupCategoryStage, unsetStage, lookupSubcategoryStage}
}

func IncrementView(id primitive.ObjectID) {
	projectCollection, ctx := database.GetCollection(collectionName)

	_, err := projectCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$inc", bson.D{{"views", 1}}},
		},
	)

	// NOTE: logging ??
	if err != nil {
		log.Fatal(err)
	}
}

func FindById(id primitive.ObjectID) (bson.M, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	lookupStage := GetLookupStage()
	matchStage := bson.D{{"$match", bson.D{{"_id", id}}}}

	stages := append(lookupStage, matchStage)

	cursor, err := projectCollection.Aggregate(ctx, stages)
	var result []bson.M
	if err != nil {
		return bson.M{}, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &result); err != nil {
		return bson.M{}, errors.NewBadRequestError(err.Error())
	}

	if len(result) == 0 {
		return bson.M{}, errors.NewNotFoundError("projectId")

	}

	return result[0], nil
}

func FindAll() ([]bson.M, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	lookupStage := GetLookupStage()

	cursor, err := projectCollection.Aggregate(ctx, lookupStage)
	var result []bson.M
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())

	}
	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	return result, nil
}

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
		{Key: "views", Value: 0},
		{Key: "createdAt", Value: time.Now()},
		{Key: "updatedAt", Value: time.Now()},
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
		"views":       0,
	}

	return result, nil
}
