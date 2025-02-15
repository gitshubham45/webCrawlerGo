package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Result represents the structure of a crawled result in MongoDB
type Result struct {
	Domain    string    `bson:"domain"`
	URLs      []string  `bson:"urls"`
	Timestamp time.Time `bson:"timestamp"`
}

// SaveResults saves the crawled results to MongoDB
func SaveResults(domain string, urls []string) {
	collection := MongoClient.Database("crawler").Collection("results")

	// Create a document to insert
	document := Result{
		Domain:    domain,
		URLs:      urls,
		Timestamp: time.Now(),
	}

	// Insert the document into the collection
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Printf("Failed to save results for domain %s to MongoDB: %v", domain, err)
		return
	}
	log.Printf("Saved results for domain %s to MongoDB", domain)
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
