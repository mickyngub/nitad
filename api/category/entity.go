package category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty`
	Title       string               `bson:"title,omitempty`
	Subcategory []primitive.ObjectID `bson:"subcategory,omitempty" json:"subcategory,omitempty"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type CategoryRequest struct {
	Title       string   `form:"title" validate:"required"`
	Subcategory []string `form:"subcategory" validate:"required"`
}
