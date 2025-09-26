package driver

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB represents the database connection
type DB struct {
	*gorm.DB
}

// NewDatabase creates a new database connection with retry mechanism
func NewDatabase() (*DB, error) {
	// Get database URL from environment variable or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default PostgreSQL connection string
		dbURL = "host=db user=user password=password dbname=go_clean_arch port=5432 sslmode=disable"
	}

	// Retry mechanism
	var db *gorm.DB
	var err error

	// Try to connect with exponential backoff
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
		if i < 4 {
			// Exponential backoff: 1s, 2s, 4s, 8s
			time.Sleep(time.Duration(1<<i) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
	}

	log.Println("Database connection established")

	return &DB{db}, nil
}

// Create implements the Database interface
func (d *DB) Create(value interface{}) error {
	result := d.DB.Create(value)
	return result.Error
}

// First implements the Database interface
func (d *DB) First(dest interface{}, conditions ...interface{}) error {
	result := d.DB.First(dest, conditions...)
	return result.Error
}

// Find implements the Database interface
func (d *DB) Find(dest interface{}, conditions ...interface{}) error {
	result := d.DB.Find(dest, conditions...)
	return result.Error
}

// Save implements the Database interface
func (d *DB) Save(value interface{}) error {
	result := d.DB.Save(value)
	return result.Error
}

// Delete implements the Database interface
func (d *DB) Delete(value interface{}, conditions ...interface{}) error {
	result := d.DB.Delete(value, conditions...)
	return result.Error
}
