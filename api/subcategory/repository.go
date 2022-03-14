package subcategory

import (
	"context"
	"fmt"
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	ListSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError)
	ListUnsetSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError)
	GetSubcategoryById(ctx context.Context, id string) (*Subcategory, errors.CustomError)
	AddSubcategory(ctx context.Context, subcate *Subcategory) (*Subcategory, errors.CustomError)
	EditSubcategory(ctx context.Context, subcate *Subcategory) (*Subcategory, errors.CustomError)
	DeleteSubcategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	InsertToCategory(ctx context.Context, subcate *Subcategory, categoryId primitive.ObjectID) (*Subcategory, errors.CustomError)
}

type subcategoryRepository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) Repository {
	return &subcategoryRepository{
		collection: client.Database(database.DatabaseName).Collection(database.COLLECTIONS["SUBCATEGORY"]),
	}
}

func (s *subcategoryRepository) ListSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError) {
	var subcates []*Subcategory
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.NewBadRequestError(err.Error())
	}

	if err = cursor.All(ctx, &subcates); err != nil {
		return nil, errors.NewBadRequestError(err.Error())
	}

	return subcates, nil
}

func (s *subcategoryRepository) ListUnsetSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError) {
	subcates := []*Subcategory{}
	pipe := mongo.Pipeline{}

	pipe = database.AppendMatchStage(pipe, "categoryId", primitive.NilObjectID)

	cursor, aggregateErr := s.collection.Aggregate(ctx, pipe)
	if aggregateErr != nil {
		return nil, errors.NewBadRequestError(aggregateErr.Error())
	}
	if curErr := cursor.All(ctx, &subcates); curErr != nil {
		return nil, errors.NewBadRequestError(curErr.Error())
	}
	return subcates, nil
}

func (s *subcategoryRepository) GetSubcategoryById(ctx context.Context, id string) (*Subcategory, errors.CustomError) {
	oid, err := database.ExtractOID(id)
	if err != nil {
		return nil, err
	}

	var subcate Subcategory
	findErr := s.collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&subcate)
	if findErr != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if findErr == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError(collectionName)
		} else {
			return nil, errors.NewBadRequestError(findErr.Error())
		}
	}
	return &subcate, nil
}

func (s *subcategoryRepository) AddSubcategory(ctx context.Context, subcate *Subcategory) (*Subcategory, errors.CustomError) {
	now := time.Now()
	subcate.CreatedAt = now
	subcate.UpdatedAt = now
	insertRes, insertErr := s.collection.InsertOne(ctx, &subcate)
	if insertErr != nil {
		return subcate, errors.NewBadRequestError(insertErr.Error())
	}

	subcate.ID = insertRes.InsertedID.(primitive.ObjectID)

	return subcate, nil
}

func (s *subcategoryRepository) EditSubcategory(ctx context.Context, subcate *Subcategory) (*Subcategory, errors.CustomError) {
	now := time.Now()
	subcate.UpdatedAt = now
	_, updateErr := s.collection.UpdateByID(
		ctx,
		subcate.ID,
		bson.D{{
			Key: "$set", Value: &subcate}})

	if updateErr != nil {
		return subcate, errors.NewBadRequestError(updateErr.Error())
	}
	return subcate, nil
}

func (s *subcategoryRepository) DeleteSubcategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	_, delErr := s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if delErr != nil {
		return errors.NewBadRequestError("Delete subcategory failed!" + delErr.Error())
	}

	return nil
}

func (s *subcategoryRepository) InsertToCategory(ctx context.Context, subcate *Subcategory, categoryId primitive.ObjectID) (*Subcategory, errors.CustomError) {
	fmt.Println("Ss", subcate.ID)
	_, updateErr := s.collection.UpdateByID(
		ctx,
		subcate.ID,
		bson.D{{
			Key: "$set", Value: bson.D{{Key: "categoryId", Value: categoryId}}}})

	if updateErr != nil {
		return subcate, errors.NewBadRequestError(updateErr.Error())
	}
	return subcate, nil
}
