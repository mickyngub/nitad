package database

import (
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExtractOID(id string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, errors.NewBadRequestError("Invalid objectId: " + id)
	}
	return objectId, nil
}
