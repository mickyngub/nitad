package spatial

import (
	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collectionName = database.COLLECTIONS["SPATIAL"]

type Spatial struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Link string             `bson:"link" json:"link"`
}

type SpatialRequest struct {
	Link string `form:"link" validate:"required"`
}
