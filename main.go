package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// URI for connect
	user := os.Getenv("MONGO_USER")
	pwd := os.Getenv("MONGO_USER_PASSWORD")
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0-yyofy.gcp.mongodb.net/test?retryWrites=true&w=majority", user, pwd)

	// Getting context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// // Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	stock := client.Database("demo").Collection("stock")

	// specify a pipeline that will only match "insert" events
	// specify the MaxAwaitTimeOption to have each attempt wait two seconds for new documents
	matchStage := bson.D{primitive.E{Key: "$match",
		Value: bson.D{primitive.E{Key: "operationType", Value: "insert"}}}}
	opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)
	changeStream, err := stock.Watch(context.TODO(), mongo.Pipeline{matchStage}, opts)

	if err != nil {
		log.Fatal(err)
	}

	for changeStream.Next(context.TODO()) {
		fmt.Println(changeStream.Current)
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
