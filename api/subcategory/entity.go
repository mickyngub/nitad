package subcategory

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subcategory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Image     string             `bson:"image" json:"image"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// type SubcategoryClean struct {
// 	ID    primitive.ObjectID `bson:"_id" json:"id"`
// 	Title string             `bson:"title" json:"title"`
// 	Image string             `bson:"image" json:"image"`
// }

type SubcategoryRequest struct {
	Title string `form:"title" validate:"required"`
	Image multipart.FileHeader
}
