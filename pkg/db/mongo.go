package db

import (
	"context"
	"fmt"

	config "github.com/ochiengotieno304/oneotp/internal/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDatabase *mongo.Database

func init() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading configs: %v", err)
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(config.MongoUri))
	if err != nil {
		panic(err)
	}

	mongoDatabase = mongoClient.Database("OneOTP")
}

func MongoClient() *mongo.Database {
	return mongoDatabase
}
