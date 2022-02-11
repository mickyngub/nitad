package functions

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
