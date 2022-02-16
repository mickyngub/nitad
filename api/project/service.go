package project

import (
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindById(oid primitive.ObjectID) (bson.M, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	lookupStage := GetLookupStage()
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: oid}}}}

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

func FindAll(oids []primitive.ObjectID) ([]bson.M, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	stages := GetLookupStage()
	for _, oid := range oids {
		stages = database.AppendMatchIdStage(stages, "subcategory._id", oid)
	}

	cursor, err := projectCollection.Aggregate(ctx, stages)
	var result []bson.M
	if err != nil {
		return []bson.M{}, errors.NewBadRequestError(err.Error())
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []bson.M{}, errors.NewBadRequestError(err.Error())
	}

	if len(result) == 0 {
		return []bson.M{}, nil
	}

	return result, nil
}

func Add(c *ProjectRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}

	subcategoryIds, categoryIds, err := ValidateAndRemoveDuplicateIds(c.Subcategory, c.Category)
	if err != nil {
		return result, err
	}

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

func Edit(oid primitive.ObjectID, upr *UpdateProjectRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}

	subcategoryIds, categoryIds, err := ValidateAndRemoveDuplicateIds(upr.Subcategory, upr.Category)
	if err != nil {
		return result, err
	}

	collection, ctx := database.GetCollection(collectionName)

	_, updateErr := collection.UpdateByID(
		ctx,
		oid,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: upr.Title},
				{Key: "description", Value: upr.Description},
				{Key: "authors", Value: upr.Authors},
				{Key: "emails", Value: upr.Emails},
				{Key: "inspiration", Value: upr.Inspiration},
				{Key: "abstract", Value: upr.Abstract},
				{Key: "images", Value: upr.Images},
				{Key: "videos", Value: upr.Videos},
				{Key: "keywords", Value: upr.Keywords},
				{Key: "category", Value: categoryIds},
				{Key: "subcategory", Value: subcategoryIds},
				{Key: "updatedAt", Value: time.Now()},
			},
		},
		})

	if updateErr != nil {
		return result, errors.NewBadRequestError("edit project error: " + updateErr.Error())
	}

	result = map[string]interface{}{
		"id":     oid,
		"title":  upr.Title,
		"images": upr.Images,
	}

	return result, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete project failed!" + err.Error())
	}

	return nil
}
