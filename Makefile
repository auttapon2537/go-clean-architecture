# Makefile for Go Clean Architecture Application

# Variables
BINARY_NAME=app
BINARY_DIR=bin
MAIN_FILE=cmd/api/main.go

# Build the application
build:
	go build -o ${BINARY_DIR}/${BINARY_NAME} ${MAIN_FILE}

# Run the application
run:
	go run ${MAIN_FILE}

# Build and run the application
build-run: build
	./${BINARY_DIR}/${BINARY_NAME}

# Clean build files
clean:
	rm -rf ${BINARY_DIR}

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Install dependencies
deps:
	go mod tidy

# Build Docker image
docker-build:
	docker build -t go-clean-architecture .

# Run Docker container
docker-run:
	docker run -p 8080:8080 go-clean-architecture

# Build and run with Docker Compose
docker-compose-up:
	docker-compose up --build

# Build and run with Docker Compose in development mode
docker-compose-dev:
	docker-compose up --build app-dev

# Stop Docker Compose
docker-compose-down:
	docker-compose down

# Run with hot reload using Air
dev:
	air -c .air.toml

# Help
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  build-run        - Build and run the application"
	@echo "  clean            - Clean build files"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  deps             - Install dependencies"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  docker-compose-up - Build and run with Docker Compose"
	@echo "  docker-compose-dev - Build and run with Docker Compose in development mode"
	@echo "  docker-compose-down - Stop Docker Compose"
	@echo "  dev              - Run with hot reload using Air"
	@echo "  help             - Show this help message"

.PHONY: build run build-run clean test test-coverage fmt vet deps docker-build docker-run docker-compose-up docker-compose-dev docker-compose-down dev help