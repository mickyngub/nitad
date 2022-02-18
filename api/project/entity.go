package project

import (
	"mime/multipart"
	"time"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID          primitive.ObjectID             `bson:"_id" json:"id"`
	Title       string                         `bson:"title" json:"title"`
	Description string                         `bson:"description" json:"description"`
	Authors     []string                       `bson:"authors" json:"authors"`
	Emails      []string                       `bson:"emails" json:"emails"`
	Inspiration string                         `bson:"inspiration" json:"inspiration"`
	Abstract    string                         `bson:"abstract" json:"abstract"`
	Images      []string                       `bson:"images,omitempty" json:"images"`
	Videos      []string                       `bson:"videos" json:"videos"`
	Keywords    []string                       `bson:"keywords" json:"keywords"`
	Status      string                         `bson:"status" json:"status"`
	Category    []category.CategoryClean       `bson:"category" json:"category"`
	Subcategory []subcategory.SubcategoryClean `bson:"subcategory" json:"subcategory"`
	Views       int                            `bson:"views" json:"views"`
	CreatedAt   time.Time                      `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time                      `bson:"updatedAt" json:"updatedAt"`
}

type ProjectRequest struct {
	Title       string                 `form:"title" validate:"required"`
	Description string                 `form:"description" validate:"required"`
	Authors     []string               `form:"authors" validate:"required"`
	Emails      []string               `form:"emails" validate:"required"`
	Inspiration string                 `form:"inspiration" validate:"required"`
	Abstract    string                 `form:"abstract" validate:"required"`
	Images      []multipart.FileHeader `form:"-"`
	Videos      []string               `form:"videos" validate:"required"`
	Keywords    []string               `form:"keywords" validate:"required"`
	Status      string                 `form:"status" validate:"required"`
	Category    []string               `form:"category" validate:"required"`
	Subcategory []string               `form:"subcategory" validate:"required"`
}

type UpdateProject struct {
	Project
	DeleteImages []string `form:"deleteImages"`
}

type UpdateProjectRequest struct {
	ProjectRequest
	DeleteImages []string `form:"deleteImages"`
}
