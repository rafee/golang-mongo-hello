package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fruit data type for inserting into collection
type Fruit struct {
	name     string
	quantity int
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

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	stock := client.Database("demo").Collection("stock")
	chgCtx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	changeStream, err := stock.Watch(chgCtx, mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(chgCtx)

	var wg sync.WaitGroup
	wg.Add(1)

	// Receive changes
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// go createChange(stock)
	go receiveChange(chgCtx, &wg, changeStream)
	wg.Wait()
}

func receiveChange(routineCtx context.Context, waitGroup *sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()

	for stream.Next(context.TODO()) {
		fmt.Println("Nothing to see here")
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
	}
}

func createChange(stock *mongo.Collection) {
	fruit := Fruit{name: "Pineapple"}
	for true {
		fruit.quantity = rand.Intn(99-10) + 10
		fmt.Println("Inserting Value", fruit)
		time.Sleep(2 * time.Second)
	}
}
