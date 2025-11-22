package mdb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB    = "hitotose"
	Games *mongo.Collection
	Users *mongo.Collection
)

func Init() {
	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// opts := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	Games = client.Database(DB).Collection("games")
	Users = client.Database(DB).Collection("users")
}
