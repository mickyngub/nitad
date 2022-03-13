package category

import (
	"time"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collectionName = database.COLLECTIONS["CATEGORY"]

type Category struct {
	ID          primitive.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string                     `bson:"title" json:"title"`
	Subcategory []*subcategory.Subcategory `bson:"subcategory" json:"subcategory"`
	CreatedAt   time.Time                  `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time                  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type CategoryDTO struct {
	ID          primitive.ObjectID `form:"-" bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `form:"title" validate:"required" bson:"title" json:"title"`
	Subcategory []string           `form:"subcategory" bson:"subcategory" json:"subcategory"`
	CreatedAt   time.Time          `form:"-" bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `form:"-" bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type CategoryLookup struct {
	ID          primitive.ObjectID   `form:"-" bson:"_id,omitempty" json:"id,omitempty"`
	Title       string               `form:"title" validate:"required" bson:"title" json:"title"`
	Subcategory []primitive.ObjectID `form:"subcategory" bson:"subcategory" json:"subcategory"`
	CreatedAt   time.Time            `form:"-" bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time            `form:"-" bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type CategorySearch struct {
	ID          primitive.ObjectID              `bson:"_id" json:"id"`
	Title       string                          `bson:"title" json:"title"`
	Subcategory []subcategory.SubcategorySearch `bson:"subcategory" json:"subcategory"`
}
