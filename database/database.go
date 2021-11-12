package database

import (
	"context"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//
var Client mongo.Client

func connetToMongo() {
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)
	clientlocal, err := mongo.Connect(context.TODO(), clientOptions)
	Client = *clientlocal
	if err != nil {
		panic(err)
	}
	err = clientlocal.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	openlog.Info("Connected to Mongodb")
}

func init() {
	connetToMongo()
}
