package project

import (
	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var SORTING = map[string]string{
	"views":     "views",
	"name":      "title",
	"updatedAt": "updatedAt",
	"createdAt": "createdAt",
}

type repositoryHelper struct{}

func (rh *repositoryHelper) AppendGetProjectStage(pipe mongo.Pipeline) mongo.Pipeline {
	pipe = database.AppendUnwindStage(pipe, "category")
	pipe = database.AppendUnwindStage(pipe, "category.subcategory")

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "subcategory"},
		{Key: "localField", Value: "category.subcategory._id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "subcategoryLookup"}}}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "category"},
		{Key: "localField", Value: "category._id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "categoryLookup"}}}})

	pipe = append(pipe, bson.D{{Key: "$project", Value: bson.D{
		{Key: "category", Value: 0},
	}}})
	return pipe
}

func (rh *repositoryHelper) AppendQueryStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	pq = rh.SetDefaultQuery(pq)

	pipe = rh.AppendSortStage(pipe, pq)
	pipe = rh.AppendPaginationStage(pipe, pq)
	return pipe
}

func (rh *repositoryHelper) AppendSortStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	return append(pipe, bson.D{
		{Key: "$sort", Value: bson.D{{
			Key: SORTING[pq.Sort], Value: pq.By,
		}}}})
}

func (rh *repositoryHelper) AppendPaginationStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	return append(pipe,
		bson.D{{Key: "$skip", Value: (pq.Page - 1) * pq.Limit}},
		bson.D{{Key: "$limit", Value: pq.Limit}})
}

func (rh *repositoryHelper) SetDefaultQuery(pq *ProjectQuery) *ProjectQuery {
	if pq.Sort == "" {
		pq.Sort = "views"
	}

	if pq.By == 0 {
		if pq.Sort == "name" {
			pq.By = 1
		} else {
			pq.By = -1
		}
	}

	if pq.Page == 0 {
		pq.Page = 1
	}

	if pq.Limit == 0 {
		pq.Limit = 15
	}
	return pq
}
