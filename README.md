# Go Clean Architecture

This is an example of a Go application implementing Clean Architecture principles. The architecture follows the separation of concerns principle, making the code more maintainable, testable, and scalable.

## Architecture Layers

### 1. Entity Layer
Contains the business objects of the application. These are the core structures that represent the domain data.

### 2. Repository Layer
Handles data access and persistence. This layer abstracts the data storage mechanism from the rest of the application.

### 3. Usecase Layer
Contains the business logic of the application. This layer orchestrates the flow of data between the repository and the handler.

### 4. Handler Layer
Handles HTTP requests and responses. This layer is responsible for parsing requests, calling usecases, and formatting responses.

### 5. Driver Layer
Contains infrastructure implementations such as database connections.

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── entity/              # Business entities
│   ├── repository/          # Data access layer
│   ├── usecase/             # Business logic layer
│   ├── handler/             # HTTP handlers
│   └── driver/              # Infrastructure implementations
├── pkg/
│   ├── utils/               # Utility functions
│   └── monitoring/          # Memory monitoring and profiling
├── go.mod                   # Go module definition
└── README.md                # This file
```

## API Endpoints

### User Management

- `POST /users/` - Create a new user
- `GET /users/:id` - Get a user by ID
- `GET /users/?email=:email` - Get a user by email
- `GET /users/all` - Get all users
- `PUT /users/:id` - Update a user
- `DELETE /users/:id` - Delete a user

## Example Requests

### Create a User
```bash
curl -X POST http://localhost:8080/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securepassword"
  }'
```

### Get a User by ID
```bash
curl http://localhost:8080/users/1
```

### Get a User by Email
```bash
curl http://localhost:8080/users/?email=john.doe@example.com
```

### Get All Users
```bash
curl http://localhost:8080/users/all
```

### Update a User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "password": "newsecurepassword"
  }'
```

### Delete a User
```bash
curl -X DELETE http://localhost:8080/users/1
```

## Development

This project includes hot reload functionality for development using Air. There are several ways to run the application in development mode:

### Using Air directly
```bash
# Install air if you haven't already
go install github.com/air-verse/air@latest

# Run with hot reload
make dev
```

### Using Docker with hot reload
```bash
# Run with Docker Compose in development mode
make docker-compose-dev
```

This will start the application with hot reload enabled. Any changes to the Go source files will automatically trigger a rebuild and restart of the application.

### Development with Docker Compose
The `docker-compose.yml` file includes two services:
- `app` - Production build of the application
- `app-dev` - Development build with hot reload

To run the development version:
```bash
docker-compose up --build app-dev
```

The development setup uses volume mapping to sync your source code with the container, allowing for hot reload functionality.

## Installation

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Run `go run cmd/api/main.go` to start the server

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - Database connection string (default: in-memory SQLite)
- `GIN_MODE` - Gin mode (default: release)
- `MONGO_URL` - MongoDB connection string (default: mongodb://admin:password@mongo:27017)

## Design Patterns Used

1. **Dependency Injection** - Dependencies are injected into each layer rather than being hardcoded
2. **Interface Abstraction** - Interfaces define contracts between layers
3. **Separation of Concerns** - Each layer has a specific responsibility
4. **SOLID Principles** - The architecture follows SOLID principles for better design

## Benefits of This Architecture

1. **Testability** - Each layer can be tested independently
2. **Maintainability** - Changes in one layer don't affect others
3. **Scalability** - Easy to add new features or modify existing ones
4. **Flexibility** - Easy to switch implementations (e.g., database, framework)

## Monitoring and Profiling

This application includes built-in memory monitoring and profiling capabilities to help detect and prevent memory leaks:

### Memory Monitoring Features

- **Real-time Memory Tracking**: Monitors memory allocation, garbage collection, and goroutine count
- **Memory Leak Detection**: Alerts when memory usage exceeds configurable thresholds
- **Per-Request Memory Tracking**: Middleware that tracks memory usage for each HTTP request
- **Health Check Endpoints**: Dedicated endpoints for monitoring application health and memory status

### Profiling Endpoints

The application includes pprof endpoints for detailed performance profiling:

- `GET /debug/pprof/` - Index of all available profiles
- `GET /debug/pprof/allocs` - Memory allocations profile
- `GET /debug/pprof/block` - Blocking operations profile
- `GET /debug/pprof/goroutine` - Goroutine profile
- `GET /debug/pprof/heap` - Heap memory profile
- `GET /debug/pprof/mutex` - Mutex contention profile
- `GET /debug/pprof/threadcreate` - Thread creation profile
- `GET /debug/pprof/cmdline` - Command line invocation
- `GET /debug/pprof/profile` - CPU profile
- `GET /debug/pprof/symbol` - Symbol lookup
- `GET /debug/pprof/trace` - Trace execution

### Health Check Endpoints

- `GET /health` - Basic health check
- `GET /health/memory` - Detailed memory usage information

### Memory Monitoring Headers

All HTTP responses include memory monitoring headers:
- `X-Memory-Before` - Memory allocation before request
- `X-Memory-After` - Memory allocation after request
- `X-Memory-Diff` - Memory allocation difference
- `X-Request-Duration` - Request processing time
- `X-Num-Goroutines` - Current number of goroutines
- `X-Goroutines-Before` - Goroutine count before request
- `X-Goroutines-After` - Goroutine count after request
- `X-Goroutines-Diff` - Goroutine count difference

## MongoDB Integration

This application now includes MongoDB integration for storing memory logs. Memory statistics are automatically captured every minute and stored in a MongoDB collection named `memory_logs` in the `go_clean_arch` database.

### MongoDB Connection

The application connects to MongoDB using the `MONGO_URL` environment variable. If not provided, it defaults to `mongodb://admin:password@mongo:27017`.

### Memory Log Structure

Memory logs stored in MongoDB have the following structure:

```json
{
  "_id": "ObjectId",
  "timestamp": "ISODate",
  "alloc": "uint64",
  "totalAlloc": "uint64",
  "sys": "uint64",
  "numGC": "uint32",
  "gcCPUFraction": "float64",
  "numGoroutine": "int"
}
```

### Docker Integration

The `docker-compose.yml` file includes a MongoDB service with the following configuration:

- Image: `mongo:6.0`
- Root username: `admin`
- Root password: `password`
- Port mapping: `27017:27017`
- Volume: `mongo_data:/data/db`