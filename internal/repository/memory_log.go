package repository

import (
	"context"
	"time"

	"github.com/example/go-clean-architecture/internal/driver"
	"github.com/example/go-clean-architecture/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MemoryLogRepository represents the repository for memory logs
type MemoryLogRepository struct {
	mongo *driver.Mongo
}

// NewMemoryLogRepository creates a new memory log repository
func NewMemoryLogRepository(mongo *driver.Mongo) *MemoryLogRepository {
	return &MemoryLogRepository{mongo: mongo}
}

// Create inserts a new memory log into MongoDB
func (r *MemoryLogRepository) Create(memoryLog *entity.MemoryLog) error {
	// Generate a new ObjectID if ID is empty
	if memoryLog.ID == "" {
		memoryLog.ID = primitive.NewObjectID().Hex()
	}

	// Set timestamp if not set
	if memoryLog.Timestamp.IsZero() {
		memoryLog.Timestamp = time.Now()
	}

	collection := r.mongo.GetCollection("go_clean_arch", "memory_logs")
	_, err := collection.InsertOne(context.Background(), memoryLog)
	return err
}

// FindByTimeRange finds memory logs within a time range
func (r *MemoryLogRepository) FindByTimeRange(start, end time.Time) ([]*entity.MemoryLog, error) {
	collection := r.mongo.GetCollection("go_clean_arch", "memory_logs")

	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var memoryLogs []*entity.MemoryLog
	if err = cursor.All(context.Background(), &memoryLogs); err != nil {
		return nil, err
	}

	return memoryLogs, nil
}

// FindAll retrieves all memory logs
func (r *MemoryLogRepository) FindAll() ([]*entity.MemoryLog, error) {
	collection := r.mongo.GetCollection("go_clean_arch", "memory_logs")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var memoryLogs []*entity.MemoryLog
	if err = cursor.All(context.Background(), &memoryLogs); err != nil {
		return nil, err
	}

	return memoryLogs, nil
}

// DeleteOlderThan deletes memory logs older than a specific time
func (r *MemoryLogRepository) DeleteOlderThan(olderThan time.Time) (int64, error) {
	collection := r.mongo.GetCollection("go_clean_arch", "memory_logs")

	filter := bson.M{
		"timestamp": bson.M{
			"$lt": olderThan,
		},
	}

	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
