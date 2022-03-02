package category

import (
	"time"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collectionName = database.COLLECTIONS["CATEGORY"]

type Category struct {
	ID          primitive.ObjectID        `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string                    `bson:"title" json:"title"`
	Subcategory []subcategory.Subcategory `bson:"subcategory" json:"subcategory"`
	CreatedAt   time.Time                 `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time                 `bson:"updatedAt" json:"updatedAt"`
}

type CategoryRequest struct {
	Title       string   `form:"title" validate:"required"`
	Subcategory []string `form:"subcategory" validate:"required"`
}

type CategorySearch struct {
	ID          primitive.ObjectID              `bson:"_id" json:"id"`
	Title       string                          `bson:"title" json:"title"`
	Subcategory []subcategory.SubcategorySearch `bson:"subcategory" json:"subcategory"`
}
