package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var client *mongo.Client

const DatabaseName = "nitad"

var COLLECTIONS = map[string]string{
	"PROJECT":     "project",
	"CATEGORY":    "category",
	"SUBCATEGORY": "subcategory",
	"ADMIN":       "admin",
	"SPATIAL":     "spatial",
}

func GetCollection(collectionName string) (*mongo.Collection, context.Context) {
	ctx := context.Background()
	collection := client.Database(DatabaseName).Collection(collectionName)

	return collection, ctx
}

func GetClient() *mongo.Client {
	return client
}

func ConnectDb(mongoURI string) *mongo.Client {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		zap.S().Fatal(err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		zap.S().Fatal(err.Error())
	}

	// List databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		zap.S().Fatal(err.Error())
	}
	zap.S().Info("databases: ", databases)

	return client
}

func DisconnectDb() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Disconnect(ctx)
	if err != nil {
		zap.S().Warn("Error disconnecting mongodb: ", err.Error())
	} else {
		zap.S().Info("disconnecting mongodb...")
	}
}
