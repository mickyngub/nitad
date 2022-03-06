package project

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/paginate"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	ListProject(ctx context.Context, pq *ProjectQuery, sids []primitive.ObjectID) ([]Project, *paginate.Paginate, errors.CustomError)

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
