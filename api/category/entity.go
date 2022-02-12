package category

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty`
	Title       string               `bson:"title,omitempty`
	Subcategory []primitive.ObjectID `bson:"subcategory,omitempty" json:"subcategory,omitempty"`
}

type CategoryRequest struct {
	ID          primitive.ObjectID
	Title       string
	Subcategory []string
}
