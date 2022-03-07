package category

import (
	"time"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAll() ([]Category, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	pipe := mongo.Pipeline{}
	pipe = database.AppendLookupStage(pipe, "subcategory")

	cursor, err := categoryCollection.Aggregate(ctx, pipe)
	var result []Category
	if err != nil {
		return result, errors.NewBadRequestError(err.Error())

	}
	if err = cursor.All(ctx, &result); err != nil {
		return result, errors.NewBadRequestError(err.Error())
	}
	// log.Println("check 6")

	return result, nil
}

func GetById(oid primitive.ObjectID) (Category, errors.CustomError) {
	categoryCollection, ctx := database.GetCollection(collectionName)

	pipe := mongo.Pipeline{}
	pipe = database.AppendMatchStage(pipe, "_id", oid)
	pipe = database.AppendLookupStage(pipe, "subcategory")

	cursor, err := categoryCollection.Aggregate(ctx, pipe)
	var result []Category
	if err != nil {
		return Category{}, errors.NewBadRequestError(err.Error())
	}
	if err = cursor.All(ctx, &result); err != nil {
		return Category{}, errors.NewBadRequestError(err.Error())
	}

	if len(result) == 0 {
		return Category{}, errors.NewNotFoundError("categoryId")

	}

	return result[0], nil
}

func Add(c *Category, sids []primitive.ObjectID) (*Category, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)
	now := time.Now()
	insertRes, insertErr := collection.InsertOne(ctx, bson.D{
		{Key: "title", Value: c.Title},
		{Key: "subcategory", Value: sids},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	})
	if insertErr != nil {
		return c, errors.NewBadRequestError(insertErr.Error())
	}

	c.ID = insertRes.InsertedID.(primitive.ObjectID)
	c.CreatedAt = now
	c.UpdatedAt = now

	return c, nil
}

func Edit(c *Category, sids []primitive.ObjectID) (*Category, errors.CustomError) {
	collection, ctx := database.GetCollection(collectionName)
	now := time.Now()
	_, updateErr := collection.UpdateByID(
		ctx,
		c.ID,
		bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "title", Value: c.Title},
				{Key: "subcategory", Value: sids},
				{Key: "updatedAt", Value: now},
			},
		},
		})

	if updateErr != nil {
		return c, errors.NewBadRequestError(updateErr.Error())
	}

	c.UpdatedAt = now

	return c, nil
}

func Delete(oid primitive.ObjectID) errors.CustomError {
	collection, ctx := database.GetCollection(collectionName)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return errors.NewBadRequestError("Delete category failed!" + err.Error())
	}

	return nil
}
