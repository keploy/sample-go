package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var msgChan chan string
var client *mongo.Client

func getTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if msgChan != nil {
		msg := time.Now().Format("15:04:05")
		msgChan <- msg

		// Store the time in MongoDB
		go func() {
			err := storeTimeInMongoDB(msg)
			if err != nil {
				fmt.Println("Error storing time in MongoDB:", err)
			}
		}()
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client connected")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgChan = make(chan string)

	defer func() {
		close(msgChan)
		msgChan = nil
		fmt.Println("Client closed connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not init http.Flusher")
	}

	for {
		select {
		case message := <-msgChan:
			fmt.Println("case message... sending message")
			fmt.Println(message)
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		}
	}

}

func storeTimeInMongoDB(msg string) error {
	collection := client.Database("your_database_name").Collection("time_collection")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, bson.M{"time": msg})
	if err != nil {
		return err
	}
	return err
}

func main() {
	router := http.NewServeMux()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	router.HandleFunc("/event", sseHandler)
	router.HandleFunc("/time", getTime)
	// fmt.Println("Server listening on :3500")
	log.Fatal(http.ListenAndServe(":3500", router))
}
