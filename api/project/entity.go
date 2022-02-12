package project

import "go.mongodb.org/mongo-driver/bson/primitive"

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
}

type ProjectRequest struct {
	Title       string
	Description string
	Authors     []string
	Emails      []string
	Inspiration string
	Abstract    string
	Images      []string
	Videos      []string
	Keywords    []string
	Category    []string
	Subcategory []string
}
