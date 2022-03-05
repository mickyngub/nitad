package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	ListCategory(ctx context.Context) ([]Category, errors.CustomError)
	// GetCateListCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	// AddCateListCategory(ctx context.Context, cate *Category) (*Category, errors.CustomError)
	// EditCateListCategory(ctx context.Context, cate *Category) (*Category, errors.CustomError)
	// DeleteCateListCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) Repository {
	return &categoryRepository{
		collection: client.Database(database.DatabaseName).Collection(database.COLLECTIONS["CATEGORY"]),
	}
}

func (c *categoryRepository) ListCategory(ctx context.Context) ([]Category, errors.CustomError) {
	pipe := mongo.Pipeline{}
	pipe = database.AppendLookupStage(pipe, "subcategory")

	cursor, err := c.collection.Aggregate(ctx, pipe)
	var cates []Category
	if err != nil {
		return cates, errors.NewBadRequestError(err.Error())

	}
	if err = cursor.All(ctx, &cates); err != nil {
		return cates, errors.NewBadRequestError(err.Error())
	}

	return cates, nil
}

func (c *categoryRepository) GetCateListCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) AddCateListCategory(ctx context.Context, cate *Category) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) EditCateListCategory(ctx context.Context, cate *Category) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) DeleteCateListCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	return nil
}
