package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Category, errors.CustomError)
	GetById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	Add(ctx context.Context, cate *Category) (*Category, errors.CustomError)
	Edit(ctx context.Context, cate *Category) (*Category, errors.CustomError)
	Delete(ctx context.Context, oid primitive.ObjectID) errors.CustomError
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) Repository {
	return &categoryRepository{
		collection: client.Database(database.DatabaseName).Collection(database.COLLECTIONS["CATEGORY"]),
	}
}

func (c *categoryRepository) FindAll(ctx context.Context) ([]Category, errors.CustomError) {
	return []Category{}, nil
}
func (c *categoryRepository) GetById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) Add(ctx context.Context, cate *Category) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) Edit(ctx context.Context, cate *Category) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) Delete(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	return nil
}
