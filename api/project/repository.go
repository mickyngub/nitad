package project

import (
	"context"
	"time"

	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	ListProject(ctx context.Context, pq *ProjectQuery, sids []primitive.ObjectID) ([]*Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx context.Context, id string) (*Project, errors.CustomError)
	AddProject(ctx context.Context, proj *Project) (*Project, errors.CustomError)
	EditProject(ctx context.Context, proj *Project) (*Project, errors.CustomError)
	DeleteProject(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	SearchProject(ctx context.Context) ([]ProjectSearch, errors.CustomError)
	IncrementView(ctx context.Context, oid primitive.ObjectID, val int)
	CountDocuments(ctx context.Context, pipe mongo.Pipeline) (int64, errors.CustomError)
}

type projectRepository struct {
	collection *mongo.Collection
	helper     *repositoryHelper
}

type Count struct {
	ID int64
}

func NewRepository(client *mongo.Client) Repository {
	return &projectRepository{
		collection: client.Database(database.DatabaseName).Collection(database.COLLECTIONS["PROJECT"]),
		helper:     &repositoryHelper{},
	}
}

func (p *projectRepository) ListProject(ctx context.Context, pq *ProjectQuery, sids []primitive.ObjectID) ([]*Project, *paginate.Paginate, errors.CustomError) {
	pipe := mongo.Pipeline{}
	projects := []*Project{}

	for _, sid := range sids {
		pipe = database.AppendMatchStage(pipe, "category.subcategory._id", sid)
	}

	count, err := p.CountDocuments(ctx, pipe)
	if err != nil {
		return projects, nil, err
	}

	queryPipe := p.helper.AppendQueryStage(pipe, pq)
	cursor, aggregateErr := p.collection.Aggregate(ctx, queryPipe)
	if aggregateErr != nil {
		return projects, nil, errors.NewBadRequestError(aggregateErr.Error())
	}
	if curErr := cursor.All(ctx, &projects); curErr != nil {
		return projects, nil, errors.NewBadRequestError(curErr.Error())
	}

	return projects, paginate.New(pq.Limit, pq.Page, count), nil
}

func (p *projectRepository) GetProjectById(ctx context.Context, id string) (*Project, errors.CustomError) {
	oid, err := database.ExtractOID(id)
	if err != nil {
		return nil, err
	}

	pipe := mongo.Pipeline{}
	pipe = database.AppendMatchStage(pipe, "_id", oid)

	cursor, mongoErr := p.collection.Aggregate(ctx, pipe)
	projects := []Project{}
	if mongoErr != nil {
		return &Project{}, errors.NewBadRequestError(mongoErr.Error())
	}
	if mongoErr = cursor.All(ctx, &projects); mongoErr != nil {
		return &Project{}, errors.NewBadRequestError(mongoErr.Error())
	}

	if len(projects) == 0 {
		return &Project{}, errors.NewNotFoundError("projectId")
	}
	return &projects[0], nil
}

func (p *projectRepository) AddProject(ctx context.Context, proj *Project) (*Project, errors.CustomError) {
	now := time.Now()
	proj.CreatedAt = now
	proj.UpdatedAt = now
	proj.Views = 0

	insertRes, insertErr := p.collection.InsertOne(ctx, proj)
	if insertErr != nil {
		return proj, errors.NewBadRequestError(insertErr.Error())
	}

	proj.ID = insertRes.InsertedID.(primitive.ObjectID)

	return proj, nil
}

func (p *projectRepository) EditProject(ctx context.Context, proj *Project) (*Project, errors.CustomError) {
	zap.S().Info("pass 8", proj.Images)

	now := time.Now()
	proj.UpdatedAt = now
	_, updateErr := p.collection.UpdateByID(
		ctx,
		proj.ID,
		bson.D{{
			Key: "$set", Value: proj}})

	if updateErr != nil {
		return proj, errors.NewBadRequestError("edit project error: " + updateErr.Error())
	}
	return proj, nil
}

func (p *projectRepository) DeleteProject(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	_, err := p.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete project failed!" + err.Error())
	}

	return nil
}

func (p *projectRepository) CountDocuments(ctx context.Context, pipe mongo.Pipeline) (int64, errors.CustomError) {
	countPipe := database.AppendCountStage(pipe)
	count := []Count{}
	cursor, aggregateErr := p.collection.Aggregate(ctx, countPipe)
	if aggregateErr != nil {
		return 0, errors.NewBadRequestError(aggregateErr.Error())
	}
	if curErr := cursor.All(ctx, &count); curErr != nil {
		return 0, errors.NewBadRequestError(curErr.Error())
	}

	if len(count) > 0 {
		return count[0].ID, nil
	} else {
		return 0, nil
	}
}

func (p *projectRepository) IncrementView(ctx context.Context, oid primitive.ObjectID, val int) {
	_, err := p.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "views", Value: val}}},
		},
	)

	if err != nil {
		zap.S().Warn("Incrementing view error: ", err.Error())
	}
}

func (p *projectRepository) SearchProject(ctx context.Context) ([]ProjectSearch, errors.CustomError) {
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
