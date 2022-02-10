package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const DatabaseName = "nitad"

var COLLECTIONS = map[string]string{
	"PROJECT":     "project",
	"CATEGORY":    "category",
	"SUBCATEGORY": "subcategory",
	"ADMIN":       "admin",
}

func GetCollection(collectionName string) (*mongo.Collection, context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database(DatabaseName).Collection(collectionName)

	return collection, ctx
}

func GetClient() *mongo.Client {
	return client
}

func ConnectDb() {
	MONGO_URI := os.Getenv("MONGO_URI")

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: fix loop bug db

	// List databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func DisconnectDb() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client.Disconnect(ctx)
}
