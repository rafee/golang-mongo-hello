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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fruit data type for inserting into collection
type Fruit struct {
	ID       primitive.ObjectID
	Name     string
	Quantity int
}

// `json:"id" bson:"_id,omitempty"`
// `json:"name" bson:"name,omitempty"`
// `json:"quantity" bson:"quantity,omitempty"`
// type ReceiveDocument struct{

// }

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	cancel()

	db := "demo"
	collec := "stock"
	stock := client.Database(db).Collection(collec)
	// Context will expire in 6 seconds. The change and receive will expire accordingly
	chgCtx, chgCancel := context.WithTimeout(context.Background(), 6*time.Second)
	// chgCancel will cancel the context before the application goes out of scope.
	// This will be activated only when the application doesn't timeout in 6 sec
	defer chgCancel()
	matchPipeline := bson.D{
		{
			"$match", bson.D{
				{"operationType", "insert"},
				{"fullDocument.quantity", bson.D{
					{"$gt", 15},
				}},
			},
		},
	}
	chgStream, err := stock.Watch(chgCtx, mongo.Pipeline{matchPipeline})
	if err != nil {
		log.Fatal(err)
	}
	defer chgStream.Close(chgCtx)

	var wg sync.WaitGroup
	wg.Add(1)

	// Receive changes
	// go createChange(stock)
	go receiveChange(chgCtx, &wg, chgStream)
	wg.Wait()
}

func receiveChange(routineCtx context.Context, waitGroup *sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()

	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			log.Fatal(err)
		}
		document := data["fullDocument"].(bson.M)
		var fruit Fruit
		bsonBytes, err := bson.Marshal(document)
		if err != nil {
			log.Fatal(err)
		}
		if err := bson.Unmarshal(bsonBytes, &fruit); err != nil {
			log.Fatal(err)
		}
		fmt.Println(fruit.Quantity)
	}

	select {
	case <-routineCtx.Done():
		fmt.Println("Operation timed out")
		return
	}
}

func createChange(stock *mongo.Collection) {
	fruit := Fruit{Name: "Pineapple"}
	for true {
		fruit.Quantity = rand.Intn(99-10) + 10
		fmt.Println("Inserting Value", fruit)
		time.Sleep(2 * time.Second)
	}
}
