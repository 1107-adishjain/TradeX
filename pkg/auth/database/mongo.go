package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

// Connect initializes the MongoDB connection pool
func Connect(uri string, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure the client options using your MONGO_URI
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	DB = client.Database(dbName)
	return DB, nil

}

// GetCollection is a clean helper function so our services can easily grab collections.
// Example: collection := database.GetCollection("users")
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
