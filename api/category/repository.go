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
	GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	AddCategory(ctx context.Context, sids []primitive.ObjectID, cate *Category) (*Category, errors.CustomError)
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

func (c *categoryRepository) GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError) {
	pipe := mongo.Pipeline{}
	pipe = database.AppendMatchStage(pipe, "_id", oid)
	pipe = database.AppendLookupStage(pipe, "subcategory")

	cursor, err := c.collection.Aggregate(ctx, pipe)
	var cates []Category
	if err != nil {
		return &Category{}, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &cates); err != nil {
		return &Category{}, errors.NewBadRequestError(err.Error())
	}

	if len(cates) == 0 {
		return &Category{}, errors.NewNotFoundError("categoryId")
	}
	return &cates[0], nil
}

func (c *categoryRepository) AddCategory(ctx context.Context, sids []primitive.ObjectID, cate *Category) (*Category, errors.CustomError) {
	insertRes, insertErr := c.collection.InsertOne(ctx, cate)
	if insertErr != nil {
		return cate, errors.NewBadRequestError(insertErr.Error())
	}
	cate.ID = insertRes.InsertedID.(primitive.ObjectID)
	return cate, nil
}

func (c *categoryRepository) EditCategory(ctx context.Context, cate *Category) (*Category, errors.CustomError) {
	return &Category{}, nil
}
func (c *categoryRepository) DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	return nil
}
