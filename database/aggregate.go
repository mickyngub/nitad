package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// $project remove all fields and return only field that in values[] with ID
func AppendProjectStage(pipe mongo.Pipeline, values []string) mongo.Pipeline {
	var finalValue bson.D
	for _, val := range values {
		finalValue = append(finalValue, bson.E{Key: val, Value: 1})
	}
	return append(pipe, bson.D{{Key: "$project", Value: finalValue}})
}

func AppendCountStage(pipe mongo.Pipeline) mongo.Pipeline {
	return append(pipe, bson.D{{Key: "$count", Value: "id"}})
}

func AppendLookupStage(pipe mongo.Pipeline, collectionName string) mongo.Pipeline {
	return append(pipe, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: collectionName},
		{Key: "localField", Value: collectionName},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: collectionName}}}})
}

// $unwind flatten array
func AppendUnwindStage(pipe mongo.Pipeline, field string) mongo.Pipeline {
	return append(pipe, bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$" + field},
		{Key: "preserveNullAndEmptyArrays", Value: true}}}})
}

func AppendUnsetStage(pipe mongo.Pipeline, field string) mongo.Pipeline {
	return append(pipe, bson.D{{Key: "$unset", Value: field}})
}

func AppendMatchStage(pipe mongo.Pipeline, field string, value interface{}) mongo.Pipeline {
	return append(pipe, bson.D{{Key: "$match", Value: bson.D{{Key: field, Value: value}}}})
}
