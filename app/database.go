package database

import (
	"context"
	// "github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
)

const CARTOON string = "test-webtoons"
const PROGRAMS string = "programs"
const SEASONS string = "seasons"

type Name struct {
	th string
	en string
}
type PROGRAM struct {
	name Name
}

var Client *mongo.Client

// Connect database.
func Connect() *mongo.Client {
	client, err := mongo.Connect(context.Background(), "mongodb://mzget:mzget1234@chitchats.ga:27017", nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	Client = client
	return client
}
