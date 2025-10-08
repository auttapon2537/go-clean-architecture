package main

import "os"

// Config holds application configuration.
type Config struct {
	port string
}

// loadConfig loads configuration from environment variables.
func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{port: port}
}
