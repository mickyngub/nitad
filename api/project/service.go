package project

import (
	"time"

	"github.com/birdglove2/nitad-backend/api/subcategory"
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

func FindAll(pq *ProjectQuery) ([]bson.M, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	subcategories, _, err := subcategory.FindByIds(pq.SubcategoryId)
	if err != nil {
		return []bson.M{}, err
	}

	stages := GetLookupStage()
	stages = AppendSortStage(stages, pq)

	for _, s := range subcategories {
		stages = database.AppendMatchIdStage(stages, "subcategory._id", s.ID)
	}

	cursor, aggregateErr := projectCollection.Aggregate(ctx, stages)
	var result []bson.M
	if aggregateErr != nil {
		return []bson.M{}, errors.NewBadRequestError(aggregateErr.Error())
	}

	if curErr := cursor.All(ctx, &result); err != nil {
		return []bson.M{}, errors.NewBadRequestError(curErr.Error())
	}

	if len(result) == 0 {
		return []bson.M{}, nil
	}

	return result, nil
}

func Add(p *Project, cid primitive.ObjectID, sids []primitive.ObjectID) (*Project, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	now := time.Now()
	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: p.Title},
		{Key: "description", Value: p.Description},
		{Key: "authors", Value: p.Authors},
		{Key: "emails", Value: p.Emails},
		{Key: "inspiration", Value: p.Inspiration},
		{Key: "abstract", Value: p.Abstract},
		{Key: "images", Value: p.Images},
		{Key: "videos", Value: p.Videos},
		{Key: "keywords", Value: p.Keywords},
		{Key: "category", Value: cid},
		{Key: "subcategory", Value: sids},
		{Key: "views", Value: 0},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	})

	if insertErr != nil {
		return p, errors.NewBadRequestError(insertErr.Error())
	}

	p.ID = insertRes.InsertedID.(primitive.ObjectID)
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Views = 0

	return p, nil
}

func Edit(oid primitive.ObjectID, upr *UpdateProjectRequest) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}

	// subcategoryIds, categoryIds, err := ValidateAndRemoveDuplicateIds(upr.Subcategory, upr.Category)
	// if err != nil {
	// 	return result, err
	// }

	// collection, ctx := database.GetCollection(collectionName)

	// _, updateErr := collection.UpdateByID(
	// 	ctx,
	// 	oid,
	// 	bson.D{{
	// 		Key: "$set", Value: bson.D{
	// 			{Key: "title", Value: upr.Title},
	// 			{Key: "description", Value: upr.Description},
	// 			{Key: "authors", Value: upr.Authors},
	// 			{Key: "emails", Value: upr.Emails},
	// 			{Key: "inspiration", Value: upr.Inspiration},
	// 			{Key: "abstract", Value: upr.Abstract},
	// 			{Key: "images", Value: upr.Images},
	// 			{Key: "videos", Value: upr.Videos},
	// 			{Key: "keywords", Value: upr.Keywords},
	// 			{Key: "category", Value: categoryIds},
	// 			{Key: "subcategory", Value: subcategoryIds},
	// 			{Key: "updatedAt", Value: time.Now()},
	// 		},
	// 	},
	// 	})

	// if updateErr != nil {
	// 	return result, errors.NewBadRequestError("edit project error: " + updateErr.Error())
	// }

	// result = map[string]interface{}{
	// 	"id":     oid,
	// 	"title":  upr.Title,
	// 	"images": upr.Images,
	// }

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
