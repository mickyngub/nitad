package project

import (
	"time"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetById(oid primitive.ObjectID) (Project, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	// pipe := GetLookupStage()
	pipe := mongo.Pipeline{}
	pipe = database.AppendMatchStage(pipe, "_id", oid)

	cursor, err := projectCollection.Aggregate(ctx, pipe)
	result := []Project{}
	if err != nil {
		return Project{}, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &result); err != nil {
		return Project{}, errors.NewBadRequestError(err.Error())
	}

	if len(result) == 0 {
		return Project{}, errors.NewNotFoundError("projectId")
	}

	return result[0], nil
}

func FindAll(pq *ProjectQuery) ([]Project, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	result := []Project{}

	_, sids, err := subcategory.FindByIds(pq.SubcategoryId)
	if err != nil {
		return result, err
	}

	// stages := GetLookupStage()
	pipe := mongo.Pipeline{}
	pipe = AppendQueryStage(pipe, pq)

	for _, sid := range sids {
		pipe = database.AppendMatchStage(pipe, "subcategory._id", sid)
	}

	cursor, aggregateErr := projectCollection.Aggregate(ctx, pipe)
	if aggregateErr != nil {
		return result, errors.NewBadRequestError(aggregateErr.Error())
	}

	if curErr := cursor.All(ctx, &result); curErr != nil {
		return result, errors.NewBadRequestError(curErr.Error())
	}

	return result, nil
}

func Add(p *Project) (*Project, errors.CustomError) {
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
		{Key: "status", Value: p.Status},
		{Key: "category", Value: p.Category},
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

func Edit(oid primitive.ObjectID, p *UpdateProject) (*UpdateProject, errors.CustomError) {

	collection, ctx := database.GetCollection(collectionName)

	now := time.Now()
	_, updateErr := collection.UpdateByID(
		ctx,
		oid,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: p.Title},
				{Key: "description", Value: p.Description},
				{Key: "authors", Value: p.Authors},
				{Key: "emails", Value: p.Emails},
				{Key: "inspiration", Value: p.Inspiration},
				{Key: "abstract", Value: p.Abstract},
				{Key: "images", Value: p.Images},
				{Key: "videos", Value: p.Videos},
				{Key: "keywords", Value: p.Keywords},
				{Key: "status", Value: p.Status},
				{Key: "category", Value: p.Category},
				{Key: "updatedAt", Value: now},
			},
		},
		})

	if updateErr != nil {
		return p, errors.NewBadRequestError("edit project error: " + updateErr.Error())
	}

	p.ID = oid
	p.UpdatedAt = now

	return p, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete project failed!" + err.Error())
	}

	return nil
}
