# Go Clean Architecture - Detailed Documentation

## Table of Contents
1. [Overview](#overview)
2. [Architecture Layers](#architecture-layers)
3. [Project Structure](#project-structure)
4. [Component Details](#component-details)
5. [Data Flow](#data-flow)
6. [API Endpoints](#api-endpoints)
7. [Database Design](#database-design)
8. [Security Considerations](#security-considerations)
9. [Deployment](#deployment)
10. [Testing](#testing)

## Overview

This project implements a clean architecture pattern in Go, following the principles of separation of concerns, dependency inversion, and testability. The architecture is divided into distinct layers, each with its own responsibility and clear boundaries.

## Architecture Layers

### 1. Entity Layer
- **Purpose**: Contains the business objects of the application
- **Responsibility**: Represent the core domain data and business rules
- **Dependencies**: None (independent)

### 2. Repository Layer
- **Purpose**: Handles data access and persistence
- **Responsibility**: Abstract the data storage mechanism from the rest of the application
- **Dependencies**: Entity layer

### 3. Usecase Layer
- **Purpose**: Contains the business logic of the application
- **Responsibility**: Orchestrate the flow of data between the repository and the handler
- **Dependencies**: Entity and Repository layers

### 4. Handler Layer
- **Purpose**: Handles HTTP requests and responses
- **Responsibility**: Parse requests, call usecases, and format responses
- **Dependencies**: Usecase layer

### 5. Driver Layer
- **Purpose**: Contains infrastructure implementations
- **Responsibility**: Provide concrete implementations for external dependencies
- **Dependencies**: External libraries

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── entity/                  # Business entities
│   │   └── user.go              # User entity and DTOs
│   ├── repository/              # Data access layer
│   │   └── user.go              # User repository interface and implementation
│   ├── usecase/                 # Business logic layer
│   │   └── user.go              # User usecase interface and implementation
│   ├── handler/                 # HTTP handlers
│   │   └── user.go              # User HTTP handlers
│   └── driver/                  # Infrastructure implementations
│       └── database.go          # Database driver implementation
├── pkg/
│   ├── utils/                   # Utility functions
│   │   └── password.go          # Password hashing utilities
│   └── monitoring/              # Memory monitoring and profiling
│       ├── memory.go            # Memory monitoring service
│       ├── pprof.go             # Pprof integration
│       └── middleware.go        # Memory monitoring middleware
├── docs/                        # Documentation
│   └── ARCHITECTURE.md          # This file
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Docker Compose configuration
├── Makefile                     # Build and run commands
├── README.md                    # Project overview
├── go.mod                       # Go module definition
└── go.sum                       # Go module checksums
```

## Component Details

### Entity Layer

The entity layer contains the core business objects. In this implementation, we have a `User` entity:

```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

We also define DTOs (Data Transfer Objects) for requests and responses:

```go
type UserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Repository Layer

The repository layer defines interfaces for data access operations:

```go
type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id uint) (*entity.User, error)
    GetByEmail(email string) (*entity.User, error)
    GetAll() ([]entity.User, error)
    Update(user *entity.User) error
    Delete(id uint) error
}
```

The implementation uses GORM for database operations:

```go
type userRepository struct {
    db Database
}

func (r *userRepository) Create(user *entity.User) error {
    return r.db.Create(user)
}
```

### Usecase Layer

The usecase layer contains business logic and orchestrates the flow of data:

```go
type UserUsecase interface {
    CreateUser(req entity.UserRequest) (*entity.UserResponse, error)
    GetUserByID(id uint) (*entity.UserResponse, error)
    GetUserByEmail(email string) (*entity.UserResponse, error)
    GetAllUsers() ([]entity.UserResponse, error)
    UpdateUser(id uint, req entity.UserRequest) (*entity.UserResponse, error)
    DeleteUser(id uint) error
}
```

Business rules are implemented here, such as password hashing:

```go
func (u *userUsecase) CreateUser(req entity.UserRequest) (*entity.UserResponse, error) {
    // Check if user already exists
    existingUser, _ := u.userRepo.GetByEmail(req.Email)
    if existingUser != nil {
        return nil, &EmailAlreadyExistsError{Email: req.Email}
    }

    // Hash the password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // Create new user entity
    user := &entity.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    // Save user to repository
    if err := u.userRepo.Create(user); err != nil {
        return nil, err
    }

    // Return response
    response := &entity.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    return response, nil
}
```

### Handler Layer

The handler layer manages HTTP requests and responses:

```go
type UserHandler struct {
    userUsecase usecase.UserUsecase
}

func (h *UserHandler) CreateHandler(c *gin.Context) {
    var req entity.UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.userUsecase.CreateUser(req)
    if err != nil {
        // Handle different error types
        switch err.(type) {
        case *usecase.EmailAlreadyExistsError:
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusCreated, response)
}
```

### Driver Layer

The driver layer provides concrete implementations for external dependencies:

```go
type DB struct {
    *gorm.DB
}

func NewDatabase() (*DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "host=db user=user password=password dbname=go_clean_arch port=5432 sslmode=disable"
    }

    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return &DB{db}, nil
}
```

## Data Flow

The data flow in this clean architecture follows a unidirectional pattern:

1. **HTTP Request** → Handler layer parses and validates the request
2. **Handler** → Calls the appropriate usecase method
3. **Usecase** → Implements business logic and calls repository methods
4. **Repository** → Interacts with the database driver
5. **Database Driver** → Performs actual database operations
6. **Response** → Flows back through the layers to the HTTP response

```
[Client] → [Handler] → [Usecase] → [Repository] → [Driver] → [Database]
                              ↖_______________________________|
```

## API Endpoints

### User Management

| Method | Endpoint     | Description          |
|--------|--------------|----------------------|
| POST   | /users/      | Create a new user    |
| GET    | /users/:id   | Get a user by ID     |
| GET    | /users/      | Get a user by email  |
| GET    | /users/all   | Get all users        |
| PUT    | /users/:id   | Update a user        |
| DELETE | /users/:id   | Delete a user        |

### Example Requests

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

#### Get a User by Email
```bash
curl http://localhost:8080/users/?email=john.doe@example.com
```

#### Get All Users
```bash
curl http://localhost:8080/users/all
```

#### Update a User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "password": "newsecurepassword"
  }'
```

#### Delete a User
```bash
curl -X DELETE http://localhost:8080/users/1
```

## Database Design

### Users Table

| Column      | Type         | Constraints              |
|-------------|--------------|--------------------------|
| id          | SERIAL       | PRIMARY KEY              |
| name        | VARCHAR(255) | NOT NULL                 |
| email       | VARCHAR(255) | UNIQUE, NOT NULL         |
| password    | VARCHAR(255) | NOT NULL                 |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP|
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP|

## Security Considerations

1. **Password Hashing**: Passwords are hashed using bcrypt before storage
2. **Input Validation**: All inputs are validated using Gin's binding features
3. **Error Handling**: Sensitive information is not exposed in error messages
4. **Environment Variables**: Configuration is managed through environment variables
5. **SQL Injection**: GORM's parameterized queries prevent SQL injection

## Deployment

### Docker Deployment

The application can be deployed using Docker:

```bash
# Build the Docker image
docker build -t go-clean-architecture .

# Run the container
docker run -p 8080:8080 go-clean-architecture
```

### Docker Compose Deployment

For a complete setup with PostgreSQL:

```bash
# Build and run with Docker Compose
docker-compose up --build
```

### Environment Variables

| Variable      | Description           | Default Value                                    |
|---------------|-----------------------|--------------------------------------------------|
| PORT          | Server port           | 8080                                             |
| DATABASE_URL  | Database connection   | host=db user=user password=password dbname=go_clean_arch port=5432 sslmode=disable |
| GIN_MODE      | Gin mode              | release                                          |

## Testing

### Unit Testing

Unit tests should be written for each layer:

1. **Entity Layer**: Test data structures and validation
2. **Repository Layer**: Test data access operations with mocks
3. **Usecase Layer**: Test business logic with mocked repositories
4. **Handler Layer**: Test HTTP handlers with mocked usecases

### Integration Testing

Integration tests should verify the interaction between layers and with the database.

### Example Test Structure

```go
func TestUserUsecase_CreateUser(t *testing.T) {
    // Setup
    mockRepo := new(mocks.UserRepository)
    usecase := NewUserUsecase(mockRepo)
    
    // Test data
    req := entity.UserRequest{
        Name:     "John Doe",
        Email:    "john.doe@example.com",
        Password: "password123",
    }
    
    // Mock expectations
    mockRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
    mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)
    
    // Execute
    result, err := usecase.CreateUser(req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, req.Name, result.Name)
    assert.Equal(t, req.Email, result.Email)
    
    // Verify mock expectations
    mockRepo.AssertExpectations(t)
}
```

## Benefits of This Architecture

1. **Testability**: Each layer can be tested independently with mocks
2. **Maintainability**: Changes in one layer don't affect others
3. **Scalability**: Easy to add new features or modify existing ones
4. **Flexibility**: Easy to switch implementations (e.g., database, framework)
5. **Separation of Concerns**: Each layer has a specific responsibility
6. **Dependency Inversion**: High-level modules don't depend on low-level modules

## Future Improvements

1. **Add more comprehensive tests**
2. **Implement logging**
3. **Add authentication and authorization**
4. **Implement caching**
5. **Add monitoring and metrics**
6. **Implement graceful shutdown**
7. **Add request validation middleware**
8. **Implement rate limiting**

## Monitoring Components

### Memory Monitoring Package

The monitoring package provides comprehensive memory tracking and profiling capabilities:

#### Memory Monitoring Service
- Tracks real-time memory allocation and garbage collection statistics
- Provides alerts when memory usage exceeds configurable thresholds
- Monitors goroutine count to detect potential leaks
- Formats memory values for human-readable output

#### Pprof Integration
- Exposes standard Go pprof endpoints for detailed performance profiling
- Includes endpoints for heap, goroutine, threadcreate, block, and mutex profiling
- Provides CPU and trace profiling capabilities

#### Monitoring Middleware
- Tracks memory usage for each HTTP request
- Monitors goroutine count before and after request processing
- Adds memory and goroutine information to HTTP response headers
- Helps identify memory-intensive endpoints

### Implementation Details

The memory monitoring is implemented as a service that can be integrated into the application's lifecycle:

```go
// Initialize memory monitor
memoryMonitor := monitoring.NewMemoryMonitor(0.8) // Alert at 80% memory usage

// Set up alert handler
memoryMonitor.SetAlertHandler(func(stats monitoring.MemoryStats) {
    log.Printf("WARN: High memory usage detected - Alloc: %s, Sys: %s",
        monitoring.FormatBytes(stats.Alloc),
        monitoring.FormatBytes(stats.Sys))
})

// Start periodic memory monitoring
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go memoryMonitor.StartMonitoring(ctx, 30*time.Second)
```

The monitoring middleware is added to the Fiber application to track per-request memory usage:

```go
// Add memory monitoring middleware
app.Use(monitoring.MemoryMiddleware(memoryMonitor))
app.Use(monitoring.SimpleGoroutineMiddleware())
```

Pprof endpoints are registered to provide detailed profiling capabilities:

```go
// Register pprof routes for profiling
monitoring.RegisterPprofRoutes(app)
```

Health check endpoints provide real-time memory status information:

```go
// Health check routes
app.Get("/health", HealthCheckHandler())
app.Get("/health/memory", monitoring.MemoryHealthCheckHandler(memoryMonitor))
```