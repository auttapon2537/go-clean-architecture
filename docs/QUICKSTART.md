# Quick Start Guide

This guide will help you get started with the Go Clean Architecture project.

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- PostgreSQL (if not using Docker)

## Project Structure

```
.
├── cmd/
│   └── api/                 # Application entry point
├── internal/                # Private application code
│   ├── entity/              # Business entities
│   ├── repository/          # Data access layer
│   ├── usecase/             # Business logic layer
│   ├── handler/             # HTTP handlers
│   └── driver/              # Infrastructure implementations
├── pkg/                     # Public libraries
│   └── utils/               # Utility functions
├── docs/                    # Documentation
├── Dockerfile               # Docker configuration
├── docker-compose.yml       # Docker Compose configuration
├── Makefile                 # Build and run commands
├── README.md                # Project overview
├── go.mod                   # Go module definition
└── go.sum                   # Go module checksums
```

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-clean-architecture
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run the Application

#### Using Go Directly

```bash
# Build the application
go build -o bin/app cmd/api/main.go

# Run the application
./bin/app
```

Or use the Makefile:

```bash
# Build and run
make build-run
```

#### Using Docker

```bash
# Build and run with Docker
make docker-build
make docker-run
```

#### Using Docker Compose (Recommended)

```bash
# Build and run with Docker Compose
make docker-compose-up
```

### 4. Test the API

Once the application is running, you can test the API endpoints:

#### Create a User

```bash
curl -X POST http://localhost:8080/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securepassword"
  }'
```

#### Get a User by ID

```bash
curl http://localhost:8080/users/1
```

#### Get All Users

```bash
curl http://localhost:8080/users/all
```

## Environment Variables

| Variable      | Description           | Default Value                                    |
|---------------|-----------------------|--------------------------------------------------|
| PORT          | Server port           | 8080                                             |
| DATABASE_URL  | Database connection   | host=db user=user password=password dbname=go_clean_arch port=5432 sslmode=disable |
| GIN_MODE      | Gin mode              | release                                          |

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Formatting

```bash
# Format code
make fmt

# Vet code
make vet
```

## Architecture Overview

This project follows the Clean Architecture pattern with four main layers:

1. **Entity Layer**: Contains business objects
2. **Repository Layer**: Handles data access
3. **Usecase Layer**: Implements business logic
4. **Handler Layer**: Manages HTTP requests/responses

Each layer depends only on the layer directly beneath it, following the Dependency Inversion Principle.

## API Endpoints

| Method | Endpoint     | Description          |
|--------|--------------|----------------------|
| POST   | /users/      | Create a new user    |
| GET    | /users/:id   | Get a user by ID     |
| GET    | /users/      | Get a user by email  |
| GET    | /users/all   | Get all users        |
| PUT    | /users/:id   | Update a user        |
| DELETE | /users/:id   | Delete a user        |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests for your changes
5. Run tests to ensure nothing is broken
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.