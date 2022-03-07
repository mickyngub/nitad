package utils

import (
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsValidObjectId(id string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return objectId, errors.NewBadRequestError("Invalid objectId")
	}
	return objectId, nil
}

func RemoveDuplicateObjectIds(ids []primitive.ObjectID) []primitive.ObjectID {
	keys := make(map[primitive.ObjectID]bool)
	list := []primitive.ObjectID{}

	for _, entry := range ids {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func RemoveDuplicateIds(ids []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range ids {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
