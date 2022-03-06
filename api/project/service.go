package project

import (
	"context"
	"time"

	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetById(oid primitive.ObjectID) (Project, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

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

type Count struct {
	ID int64
}

func (p *projectRepository) SearchAll(ctx context.Context) ([]ProjectSearch, errors.CustomError) {
	var projs []ProjectSearch
	pipe := mongo.Pipeline{}
	pipe = database.AppendProjectStage(pipe, []string{"title"})

	cursor, err := p.collection.Aggregate(ctx, pipe)
	if err != nil {
		return projs, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &projs); err != nil {
		return projs, errors.NewBadRequestError(err.Error())
	}

	return projs, nil
}

func FindAll(pq *ProjectQuery) ([]Project, paginate.Paginate, errors.CustomError) {
	projectCollection, ctx := database.GetCollection(collectionName)

	pagin := paginate.Paginate{}
	result := []Project{}

	_, sids, err := subcategory.FindByIds(pq.SubcategoryId)
	if err != nil {
		return result, pagin, err
	}

	pipe := mongo.Pipeline{}

	for _, sid := range sids {
		pipe = database.AppendMatchStage(pipe, "category.subcategory._id", sid)
	}

	countPipe := database.AppendCountStage(pipe)

	count := []Count{}
	cursor, aggregateErr := projectCollection.Aggregate(ctx, countPipe)
	if aggregateErr != nil {
		return result, pagin, errors.NewBadRequestError(aggregateErr.Error())
	}

	if curErr := cursor.All(ctx, &count); curErr != nil {
		return result, pagin, errors.NewBadRequestError(curErr.Error())
	}

	queryPipe := AppendQueryStage(pipe, pq)
	cursor, aggregateErr = projectCollection.Aggregate(ctx, queryPipe)
	if aggregateErr != nil {
		return result, pagin, errors.NewBadRequestError(aggregateErr.Error())
	}

	if curErr := cursor.All(ctx, &result); curErr != nil {
		return result, pagin, errors.NewBadRequestError(curErr.Error())
	}

	if len(count) > 0 {
		pagin = *(paginate.New(pq.Limit, pq.Page, count[0].ID))
	} else {
		pagin = *(paginate.New(pq.Limit, pq.Page, 0))
	}

	return result, pagin, nil
}

func Add(p *Project) (*Project, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	now := time.Now()

	p.Views = 0
	p.CreatedAt = now
	p.UpdatedAt = now
	insertRes, insertErr := collection.InsertOne(ctx, &p)

	if insertErr != nil {
		return p, errors.NewBadRequestError(insertErr.Error())
	}

	p.ID = insertRes.InsertedID.(primitive.ObjectID)
	return p, nil
}

func Edit(p *Project) (*Project, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)

	now := time.Now()
	p.UpdatedAt = now
	_, updateErr := collection.UpdateByID(ctx, p.ID, bson.D{{Key: "$set", Value: &p}})

	if updateErr != nil {
		return p, errors.NewBadRequestError("edit project error: " + updateErr.Error())
	}

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
