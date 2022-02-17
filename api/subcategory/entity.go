package subcategory

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subcategory struct {
	ID        primitive.ObjectID `bson:"_id",omitempty`
	Title     string             `bson:"title" validate:"required"`
	Image     string             `bson:"image"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type SubcategoryRequest struct {
	Title string `form:"title" validate:"required"`
	Image multipart.FileHeader
}
