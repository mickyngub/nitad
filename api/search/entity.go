package search

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/project"
)

type Search struct {
	Project  []project.ProjectSearch   `bson:"project" json:"project"`
	Category []category.CategorySearch `bson:"category" json:"category"`
}
