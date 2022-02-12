package subcategory

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subcategory struct {
	ID    primitive.ObjectID `bson:"_id,omitempty`
	Title string             `bson:"title,omitempty`
	Image string             `bson:"image,omitempty`
}
