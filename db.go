package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout           = 5
	connectionStringTemplate = "mongodb+srv://%s:%s@%s"
)

// GetConnection Retrieves a client to the MongoDB
func getConnection() (*mongo.Client, context.Context, context.CancelFunc) {

	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	log.Print(connectionURI)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client, ctx, cancel
}

// GetAllTracks Retrives all tracks from the db
func GetAllTracks() ([]*Track, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	var tracks []*Track

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	cursor, err := client.Database(databaseENV).Collection("track").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &tracks)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return tracks, nil
}

// GetOneTrack - search for specif track in the db
func GetOneTrack(name string) (bson.M, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	var track bson.M

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	log.Print(name)

	err := client.Database(databaseENV).Collection("track").FindOne(ctx, bson.M{"name": name}).Decode(&track)
	if err != nil {
		log.Print(err)
	}

	return track, nil
}

//Create creating a track in a mongo
func Create(track *Track) (primitive.ObjectID, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	track._ID = primitive.NewObjectID()

	result, err := client.Database(databaseENV).Collection("track").InsertOne(ctx, track)
	println(result)
	if err != nil {
		log.Printf("Could not create Track: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

// GetCredentials - checks the credetials
func GetCredentials(cred *Credentials) (bson.M, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	var result bson.M

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(databaseENV)
	collection := db.Collection("user")
	err := collection.FindOne(ctx, bson.M{"username": cred.Username}).Decode(&result)
	if err != nil {
		log.Print(err)
	}

	return result, nil
}

// GetAllEvents Retrives all events from the db
func GetAllEvents() ([]*Track, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	var tracks []*Track

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	cursor, err := client.Database(databaseENV).Collection("event").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &tracks)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return tracks, nil
}

// GetOneEvent - search for specif event in the db
func GetOneEvent(name string) (bson.M, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	var event bson.M

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	log.Print(name)

	err := client.Database(databaseENV).Collection("event").FindOne(ctx, bson.M{"name": name}).Decode(&event)
	if err != nil {
		log.Print(err)
	}

	return event, nil
}

// CreateEvent creating a track in a mongo
func CreateEvent(event *Event) (primitive.ObjectID, error) {
	databaseENV := os.Getenv("MONGODB_DATABASE")
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	event._ID = primitive.NewObjectID()

	result, err := client.Database(databaseENV).Collection("event").InsertOne(ctx, event)
	println(result)
	if err != nil {
		log.Printf("Could not create Track: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}
