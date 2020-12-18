package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var AppContext = context.TODO()

// collections
var BookCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@cluster0.zmz84.mongodb.net/test?authSource=admin&replicaSet=atlas-8x6aaq-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true")
	client, err := mongo.Connect(AppContext, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(AppContext, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You connected to your mongo database.")

	BookCollection = client.Database("test").Collection("books")
}
