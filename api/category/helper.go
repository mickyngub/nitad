package category

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
)

// receive array of categoryIds, then
// find and return non-duplicated categories, and their ids
func FindById(cid string) (CategoryProject, errors.CustomError) {
	var categoryProject CategoryProject

	oid, err := functions.IsValidObjectId(cid)
	if err != nil {
		return categoryProject, err
	}

	bson, err := database.FindById(oid, collectionName)
	category := BsonToCategory(bson)
	if err != nil {
		return categoryProject, err
	}
	categoryProject.ID = category.ID
	categoryProject.Title = category.Title

	return categoryProject, nil
}

// validate requested string of a single categoryId
// and return valid objectId, otherwise error
func ValidateId(cid string) (Category, errors.CustomError) {
	var c Category
	objectId, err := functions.IsValidObjectId(cid)
	if err != nil {
		return c, err
	}

	if c, err = GetById(objectId); err != nil {
		return c, err
	}

	return c, nil
}

// convert bson to category
func BsonToCategory(b bson.M) Category {
	var s Category
	bsonBytes, _ := bson.Marshal(b)
	bson.Unmarshal(bsonBytes, &s)
	return s
}
