package category

import (
	"context"
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	ListCategory(ctx context.Context) ([]Category, errors.CustomError)
	GetCategoryById(ctx context.Context, oid primitive.ObjectID) (*Category, errors.CustomError)
	AddCategory(ctx context.Context, cate *CategoryDTO, osids []primitive.ObjectID) (*CategoryDTO, errors.CustomError)
	EditCategory(ctx context.Context, cate *CategoryDTO, osids []primitive.ObjectID) (*CategoryDTO, errors.CustomError)
	DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	SearchCategory(ctx context.Context) ([]CategorySearch, errors.CustomError)
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

func (c *categoryRepository) AddCategory(ctx context.Context, cate *CategoryDTO, osids []primitive.ObjectID) (*CategoryDTO, errors.CustomError) {
	now := time.Now()
	insertRes, insertErr := c.collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: cate.Title},
		{Key: "subcategory", Value: osids},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	})
	if insertErr != nil {
		return cate, errors.NewBadRequestError(insertErr.Error())
	}

	cate.ID = insertRes.InsertedID.(primitive.ObjectID)
	return cate, nil
}

func (c *categoryRepository) EditCategory(ctx context.Context, cate *CategoryDTO, osids []primitive.ObjectID) (*CategoryDTO, errors.CustomError) {
	_, updateErr := c.collection.UpdateByID(
		ctx,
		cate.ID,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: cate.Title},
				{Key: "subcategory", Value: osids},
				{Key: "updatedAt", Value: time.Now()},
			},
		},
		})

	if updateErr != nil {
		return cate, errors.NewBadRequestError(updateErr.Error())
	}

	return cate, nil
}

func (c *categoryRepository) DeleteCategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	_, err := c.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete category failed!" + err.Error())
	}
	return nil
}

func (c *categoryRepository) SearchCategory(ctx context.Context) ([]CategorySearch, errors.CustomError) {
	var result []CategorySearch
	pipe := mongo.Pipeline{}
	pipe = database.AppendLookupStage(pipe, "subcategory")
	pipe = database.AppendProjectStage(pipe, []string{"title", "subcategory"})

	cursor, err := c.collection.Aggregate(ctx, pipe)
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}

	return result, nil
}
