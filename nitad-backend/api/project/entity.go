package project

import (
	"mime/multipart"
	"time"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collectionName = database.COLLECTIONS["PROJECT"]

type Project struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string              `bson:"title" json:"title"`
	Description  string              `bson:"description" json:"description"`
	Authors      []string            `bson:"authors" json:"authors"`
	Emails       []string            `bson:"emails" json:"emails"`
	Inspiration  string              `bson:"inspiration" json:"inspiration"`
	Abstract     string              `bson:"abstract" json:"abstract"`
	Images       []string            `bson:"images,omitempty" json:"images"`
	Videos       []string            `bson:"videos" json:"videos"`
	Keywords     []string            `bson:"keywords" json:"keywords"`
	Report       string              `bson:"report" json:"report"`
	VirtualLink  string              `bson:"virtualLink" json:"virtualLink"`
	Status       string              `bson:"status" json:"status"`
	Category     []category.Category `bson:"category" json:"category"`
	Views        int                 `bson:"views" json:"views"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeleteImages []string            `bson:"deleteImages,omitempty" json:"deleteImages,omitempty"`
}

type ProjectDTO struct {
	ID           primitive.ObjectID      `form:"-"`
	Title        string                  `form:"title" validate:"required"`
	Description  string                  `form:"description" validate:"required"`
	Authors      []string                `form:"authors" validate:"required"`
	Emails       []string                `form:"emails" validate:"required"`
	Inspiration  string                  `form:"inspiration" validate:"required"`
	Abstract     string                  `form:"abstract" validate:"required"`
	Images       []*multipart.FileHeader `form:"-"`
	Videos       []string                `form:"videos" validate:"required"`
	Keywords     []string                `form:"keywords"`
	Report       *multipart.FileHeader   `form:"-"`
	VirtualLink  string                  `form:"virtualLink"`
	Status       string                  `form:"status" validate:"required"`
	Category     []string                `form:"category" validate:"required"`
	Subcategory  []string                `form:"subcategory" validate:"required"`
	DeleteImages []string                `form:"deleteImages"`
}

type ProjectLookup struct {
	ID                primitive.ObjectID        `bson:"_id,omitempty" json:"id,omitempty"`
	Title             string                    `bson:"title" json:"title"`
	Description       string                    `bson:"description" json:"description"`
	Authors           []string                  `bson:"authors" json:"authors"`
	Emails            []string                  `bson:"emails" json:"emails"`
	Inspiration       string                    `bson:"inspiration" json:"inspiration"`
	Abstract          string                    `bson:"abstract" json:"abstract"`
	Images            []string                  `bson:"images,omitempty" json:"images"`
	Videos            []string                  `bson:"videos" json:"videos"`
	Keywords          []string                  `bson:"keywords" json:"keywords"`
	Report            string                    `bson:"report" json:"report"`
	VirtualLink       string                    `bson:"virtualLink" json:"virtualLink"`
	Status            string                    `bson:"status" json:"status"`
	Views             int                       `bson:"views" json:"views"`
	CreatedAt         time.Time                 `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time                 `bson:"updatedAt" json:"updatedAt"`
	DeleteImages      []string                  `bson:"deleteImages,omitempty" json:"deleteImages,omitempty"`
	CategoryLookup    []category.CategoryLookup `bson:"categoryLookup,omitempty" json:"categoryLookup,omitempty"`
	SubcategoryLookup []subcategory.Subcategory `bson:"subcategoryLookup,omitempty" json:"subcategoryLookup,omitempty"`
}

type ProjectSearch struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `bson:"title" json:"title"`
}

type ProjectQuery struct {
	SubcategoryId []string `query:"subcategoryId"`
	Sort          string   `query:"sort"`
	By            int      `query:"by"`
	Page          int      `query:"page"`
	Limit         int      `query:"limit"`
}
