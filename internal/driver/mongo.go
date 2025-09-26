package driver

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo represents the MongoDB connection
type Mongo struct {
	Client *mongo.Client
}

// NewMongo creates a new MongoDB connection with retry mechanism
func NewMongo() (*Mongo, error) {
	// Get MongoDB URL from environment variable or use default
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		// Default MongoDB connection string with authentication
		mongoURL = "mongodb://admin:password@mongo:27017"
	}

	var client *mongo.Client
	var err error

	// Retry mechanism
	for i := 0; i < 5; i++ {
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoURL)

		// Set connection timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Connect to MongoDB
		client, err = mongo.Connect(ctx, clientOptions)
		cancel()

		if err == nil {
			// Check the connection
			ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
			err = client.Ping(ctx, nil)
			cancel()

			if err == nil {
				break
			}
		}

		log.Printf("Failed to connect to MongoDB (attempt %d): %v", i+1, err)
		if i < 4 {
			// Exponential backoff: 1s, 2s, 4s, 8s
			time.Sleep(time.Duration(1<<i) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB after 5 attempts: %w", err)
	}

	log.Println("MongoDB connection established")

	return &Mongo{Client: client}, nil
}

// GetCollection returns a MongoDB collection
func (m *Mongo) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}

// Close closes the MongoDB connection
func (m *Mongo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return m.Client.Disconnect(ctx)
}
