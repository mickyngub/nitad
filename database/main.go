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


COLLECTIONS := map[string]string{
	"PROJECT":     "project",
	"CATEGORY":    "category",
	"SUBCATEGORY": "subcategory",
	"ADMIN":       "admin",
}

func connectDb() (*mongo.Client, []string) {
	MONGO_URI := os.Getenv("MONGO_URI")

	fmt.Println(COLLECTIONS)
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	return client, databases
}

