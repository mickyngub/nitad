package project

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty`
	Title       string               `bson:"title,omitempty`
	Description string               `bson:"description,omitempty`
	Authors     []string             `bson:"authors,omitempty`
	Emails      []string             `bson:"emails,omitempty`
	Inspiration string               `bson:"inspiration,omitempty`
	Abstract    string               `bson:"abstract,omitempty`
	Images      []string             `bson:"images,omitempty`
	Videos      []string             `bson:"videos,omitempty`
	Keywords    []string             `bson:"keywords,omitempty`
	Category    []primitive.ObjectID `bson:"category,omitempty" json:"category,omitempty"`
	Subcategory []primitive.ObjectID `bson:"subcategory,omitempty" json:"subcategory,omitempty"`
	Views       int                  `bson:"views,omitempty" json:"views,omitempty"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type ProjectRequest struct {
	Title       string   `form:"title" validate:"required"`
	Description string   `form:"description" validate:"required"`
	Authors     []string `form:"authors" validate:"required"`
	Emails      []string `form:"emails" validate:"required"`
	Inspiration string   `form:"inspiration" validate:"required"`
	Abstract    string   `form:"abstract" validate:"required"`
	Images      []string `form:"images" validate:"required"`
	Videos      []string `form:"videos" validate:"required"`
	Keywords    []string `form:"keywords" validate:"required"`
	Category    []string `form:"category" validate:"required"`
	Subcategory []string `form:"subcategory" validate:"required"`
}

type UpdateProjectRequest struct {
	ProjectRequest
	DeleteImages []string `form:deleteImages`
}
