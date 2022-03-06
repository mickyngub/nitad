package project

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	ListProject(ctx context.Context, pq *ProjectQuery, sids []primitive.ObjectID) ([]Project, *paginate.Paginate, errors.CustomError)
	GetProjectById(ctx context.Context, oid primitive.ObjectID) (*Project, errors.CustomError)

	IncrementView(ctx context.Context, oid primitive.ObjectID, val int)
	CountDocuments(ctx context.Context, pipe mongo.Pipeline) (int64, errors.CustomError)
}

type projectRepository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) Repository {
	return &projectRepository{
		collection: client.Database(database.DatabaseName).Collection(database.COLLECTIONS["PROJECT"]),
	}
}

func (p *projectRepository) ListProject(ctx context.Context, pq *ProjectQuery, sids []primitive.ObjectID) ([]Project, *paginate.Paginate, errors.CustomError) {
	pipe := mongo.Pipeline{}
	projects := []Project{}

	for _, sid := range sids {
		pipe = database.AppendMatchStage(pipe, "category.subcategory._id", sid)
	}

	count, err := p.CountDocuments(ctx, pipe)
	if err != nil {
		return projects, nil, err
	}

	queryPipe := AppendQueryStage(pipe, pq)
	cursor, aggregateErr := p.collection.Aggregate(ctx, queryPipe)
	if aggregateErr != nil {
		return projects, nil, errors.NewBadRequestError(aggregateErr.Error())
	}
	if curErr := cursor.All(ctx, &projects); curErr != nil {
		return projects, nil, errors.NewBadRequestError(curErr.Error())
	}

	return projects, paginate.New(pq.Limit, pq.Page, count), nil
}

func (p *projectRepository) GetProjectById(ctx context.Context, oid primitive.ObjectID) (*Project, errors.CustomError) {
	pipe := mongo.Pipeline{}
	pipe = database.AppendMatchStage(pipe, "_id", oid)

	cursor, err := p.collection.Aggregate(ctx, pipe)
	projects := []Project{}
	if err != nil {
		return &Project{}, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &projects); err != nil {
		return &Project{}, errors.NewBadRequestError(err.Error())
	}

	if len(projects) == 0 {
		return &Project{}, errors.NewNotFoundError("projectId")
	}
	return &projects[0], nil
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
