package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

//Connect used to connect to the mongodb database
func Connect(database string) *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://" + os.Getenv("MONGODB_HOST") + ":" + os.Getenv("MONGODB_PORT") + "/?readPreference=primary&appname=MongoDB%20Compass&ssl=false")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongodb database!")
	return client.Database(database)
}
