package category

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty`
	Title       string               `bson:"title,omitempty`
	Images      []string             `bson:"images,omitempty`
	Subcategory []primitive.ObjectID `bson:"subcategory,omitempty" json:"subcategory,omitempty"`
}
