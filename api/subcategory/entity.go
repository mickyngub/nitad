package subcategory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subcategory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty`
	Title     string             `bson:"title,omitempty`
	Image     string             `bson:"image,omitempty`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	updatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
