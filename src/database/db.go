package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
var db *mongo.Database

const (
	hostname       = "mongodb://localhost:27017"
	dbURI          = "mongodb://localhost:27017"
	dbName         = "demo_todo"
	collectionName = "todo"
	todoCollection = "todo"
)

func ConnectData() *mongo.Client {
	c, error := mongo.Connect(options.Client().ApplyURI(dbURI))
	checkErr(error, "failed connect to mongodb")
	client = c
	db = client.Database(dbName)

	return client
}

func DisconnectData() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Println("Database client disconnected")
		panic(err)
	}
}

func TodoCollection() *mongo.Collection {
	collection := db.Collection(todoCollection)
	return collection
}

func checkErr(e error, customMsg string) {
	if e != nil {
		fmt.Println(customMsg)
		log.Fatal(e)
	}
}
