package category

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO: reuse these 2 funcs with subcategory's validating func
// validate requested string of categoryIds
// and return valid []objectId, otherwise error
func ValidateIds(cids []string) ([]primitive.ObjectID, errors.CustomError) {
	objectIds := make([]primitive.ObjectID, len(cids))

	for i, cid := range cids {
		objectId, err := ValidateId(cid)
		if err != nil {
			return objectIds, err
		}

		objectIds[i] = objectId
	}

	return objectIds, nil
}

// validate requested string of a single categoryId
// and return valid objectId, otherwise error
func ValidateId(cid string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := functions.IsValidObjectId(cid)
	if err != nil {
		return objectId, err
	}

	// if err != nil >> id is not found
	if _, err = FindById(objectId); err != nil {
		return objectId, err
	}

	return objectId, nil
}

func BsonToCategory(b bson.M) CategoryRequest {
	// convert bson to subcategory
	var s CategoryRequest
	bsonBytes, _ := bson.Marshal(b)
	bson.Unmarshal(bsonBytes, &s)
	return s
}
