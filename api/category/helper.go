package category

import (
	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// use for checking ids from ProjectRequest
// receive array of categoryIds, then
// find and return non-duplicated categories, and their ids
// return []CategoryClean
func FindByIds(cids []string) ([]CategoryClean, []primitive.ObjectID, errors.CustomError) {
	var objectIds []primitive.ObjectID
	var categories []CategoryClean

	cids = functions.RemoveDuplicateIds(cids)

	for _, cid := range cids {
		oid, err := functions.IsValidObjectId(cid)
		if err != nil {
			return categories, objectIds, err
		}

		bson, err := database.GetElementById(oid, collectionName)
		category := BsonToCategory(bson)
		if err != nil {
			return categories, objectIds, err
		}

		objectIds = append(objectIds, oid)
		categories = append(categories, CategoryClean{
			ID:    category.ID,
			Title: category.Title,
		})

	}

	return categories, objectIds, nil

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
