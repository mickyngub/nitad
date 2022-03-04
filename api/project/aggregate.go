package project

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var SORTING = map[string]string{
	"views":     "views",
	"name":      "title",
	"updatedAt": "updatedAt",
	"createdAt": "createdAt",
}

func AppendQueryStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	pq = SetDefaultQuery(pq)

	pipe = AppendSortStage(pipe, pq)
	pipe = AppendPaginationStage(pipe, pq)
	return pipe
}

func AppendSortStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	return append(pipe, bson.D{
		{Key: "$sort", Value: bson.D{{
			Key: SORTING[pq.Sort], Value: pq.By,
		}}}})
}

func AppendPaginationStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	return append(pipe,
		bson.D{{Key: "$skip", Value: (pq.Page - 1) * pq.Limit}},
		bson.D{{Key: "$limit", Value: pq.Limit}})
}

func SetDefaultQuery(pq *ProjectQuery) *ProjectQuery {
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
