package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Result represents the structure of a crawled result in MongoDB
type Result struct {
	Domain    string    `bson:"domain"`
	URLs      []string  `bson:"urls"`
	Timestamp time.Time `bson:"timestamp"`
}

// SaveResults upserts the crawled results in MongoDB
func SaveResults(domain string, urls []string) {
	collection := MongoClient.Database("crawler").Collection("results")

	if urls == nil {
		urls = []string{} // Initialize as an empty slice
	}

	// Create an update query to add new URLs to the existing array
	filter := bson.M{"domain": domain}
	update := bson.M{
		"$addToSet":    bson.M{"urls": bson.M{"$each": urls}}, // Add unique URLs to the array
		"$setOnInsert": bson.M{"timestamp": time.Now()},       // Set timestamp only on insert
	}
	opts := options.Update().SetUpsert(true)

	// Perform the upsert operation
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Failed to save results for domain %s to MongoDB: %v", domain, err)
		return
	}
	log.Printf("Upserted results for domain %s in MongoDB", domain)
}

// GetResults retrieves all results from MongoDB
func GetResults() ([]Result, error) {
	collection := MongoClient.Database("crawler").Collection("results")

	// Query all documents
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []Result
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
