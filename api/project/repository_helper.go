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

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "subcategory"},
		{Key: "localField", Value: "category.subcategory._id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "category.subcategory"}}}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "category"},
		{Key: "localField", Value: "category._id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "categoryLookup"}}}})

	pipe = database.AppendUnwindStage(pipe, "categoryLookup")

	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.D{
		{Key: "category.title", Value: "$categoryLookup.title"},
		{Key: "category.createdAt", Value: "$categoryLookup.createdAt"},
		{Key: "category.updatedAt", Value: "$categoryLookup.updatedAt"},
	}}})

	pipe = append(pipe, bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$_id"},
		{Key: "title", Value: bson.D{{Key: "$first", Value: "$title"}}},
		{Key: "description", Value: bson.D{{Key: "$first", Value: "$description"}}},
		{Key: "authors", Value: bson.D{{Key: "$first", Value: "$authors"}}},
		{Key: "emails", Value: bson.D{{Key: "$first", Value: "$emails"}}},
		{Key: "inspiration", Value: bson.D{{Key: "$first", Value: "$inspiration"}}},
		{Key: "abstract", Value: bson.D{{Key: "$first", Value: "$abstract"}}},
		{Key: "images", Value: bson.D{{Key: "$first", Value: "$images"}}},
		{Key: "videos", Value: bson.D{{Key: "$first", Value: "$videos"}}},
		{Key: "keywords", Value: bson.D{{Key: "$first", Value: "$keywords"}}},
		{Key: "report", Value: bson.D{{Key: "$first", Value: "$report"}}},
		{Key: "virtualLink", Value: bson.D{{Key: "$first", Value: "$virtualLink"}}},
		{Key: "status", Value: bson.D{{Key: "$first", Value: "$status"}}},
		{Key: "views", Value: bson.D{{Key: "$first", Value: "$views"}}},
		{Key: "createdAt", Value: bson.D{{Key: "$first", Value: "$createdAt"}}},
		{Key: "updatedAt", Value: bson.D{{Key: "$first", Value: "$updatedAt"}}},
		{Key: "category", Value: bson.D{{Key: "$push", Value: "$category"}}},
	}}})

	return pipe
}

func (rh *repositoryHelper) AppendGetProjectStageOld(pipe mongo.Pipeline) mongo.Pipeline {
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
