package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/go-clean-architecture/internal/driver"
	"github.com/example/go-clean-architecture/internal/entity"
	"github.com/example/go-clean-architecture/internal/handler"
	"github.com/example/go-clean-architecture/internal/repository"
	"github.com/example/go-clean-architecture/internal/usecase"
	"github.com/example/go-clean-architecture/pkg/monitoring"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// App represents the application with all its components
type App struct {
	fiberApp      *fiber.App
	db            *driver.DB
	mongo         *driver.Mongo
	memoryMonitor *monitoring.MemoryMonitor
	memoryLogRepo *repository.MemoryLogRepository
	userRepo      repository.UserRepository
	userUsecase   usecase.UserUsecase
	userHandler   *handler.UserHandler
	ctx           context.Context
	cancel        context.CancelFunc
}

// Config holds application configuration
type Config struct {
	port string
}

func main() {
	// Load configuration
	config := loadConfig()

	// Initialize application
	app, err := initializeApp()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer app.cleanup()

	// Start memory monitoring
	app.startMemoryMonitoring()

	// Start memory logging
	app.startMemoryLogging()

	// Setup routes
	app.setupRoutes()

	// Print registered routes for debugging
	app.printRoutes()

	// Start server
	if err := app.startServer(config.port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// Wait for interrupt signal for graceful shutdown
	app.waitForShutdown()
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		port: port,
	}
}

// initializeApp initializes all application components
func initializeApp() (*App, error) {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize memory monitor
	memoryMonitor := monitoring.NewMemoryMonitor(0.8) // Alert at 80% memory usage

	// Set up alert handler
	memoryMonitor.SetAlertHandler(func(stats monitoring.MemoryStats) {
		log.Printf("WARN: High memory usage detected - Alloc: %s, Sys: %s",
			monitoring.FormatBytes(stats.Alloc),
			monitoring.FormatBytes(stats.Sys))
	})

	// Initialize PostgreSQL database
	db, err := driver.NewDatabase()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run auto migration
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize MongoDB
	mongo, err := driver.NewMongo()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Initialize repositories
	memoryLogRepo := repository.NewMemoryLogRepository(mongo)
	userRepo := repository.NewUserRepository(db)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)

	// Initialize Fiber app
	fiberApp := fiber.New()

	// Add middleware
	fiberApp.Use(logger.New())
	fiberApp.Use(monitoring.MemoryMiddleware(memoryMonitor))
	fiberApp.Use(monitoring.SimpleGoroutineMiddleware())

	// Register pprof routes for profiling
	monitoring.RegisterPprofRoutes(fiberApp)

	return &App{
		fiberApp:      fiberApp,
		db:            db,
		mongo:         mongo,
		memoryMonitor: memoryMonitor,
		memoryLogRepo: memoryLogRepo,
		userRepo:      userRepo,
		userUsecase:   userUsecase,
		userHandler:   userHandler,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// startMemoryMonitoring starts the periodic memory monitoring
func (app *App) startMemoryMonitoring() {
	go app.memoryMonitor.StartMonitoring(app.ctx, 30*time.Second)
}

// startMemoryLogging starts periodic memory logging to MongoDB
func (app *App) startMemoryLogging() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				stats := app.memoryMonitor.GetMemoryStats()
				log.Printf("MEMORY STATS - Alloc: %s, TotalAlloc: %s, Sys: %s, NumGC: %d, GCCPUFraction: %.4f, NumGoroutine: %d",
					monitoring.FormatBytes(stats.Alloc),
					monitoring.FormatBytes(stats.TotalAlloc),
					monitoring.FormatBytes(stats.Sys),
					stats.NumGC,
					stats.GCCPUFraction,
					stats.NumGoroutine)

				// Store memory stats in MongoDB
				memoryLog := &entity.MemoryLog{
					Alloc:         stats.Alloc,
					TotalAlloc:    stats.TotalAlloc,
					Sys:           stats.Sys,
					NumGC:         stats.NumGC,
					GCCPUFraction: stats.GCCPUFraction,
					NumGoroutine:  stats.NumGoroutine,
				}

				if err := app.memoryLogRepo.Create(memoryLog); err != nil {
					log.Printf("ERROR: Failed to store memory log in MongoDB: %v", err)
				}
			}
		}
	}()
}

// setupRoutes sets up all application routes
func (app *App) setupRoutes() {
	// Health check routes
	app.fiberApp.Get("/health", HealthCheckHandler())
	app.fiberApp.Get("/health/memory", monitoring.MemoryHealthCheckHandler(app.memoryMonitor))

	// Setup user routes
	setupUserRoutes(app.fiberApp, app.userHandler)
}

// setupUserRoutes sets up user-related routes
func setupUserRoutes(router *fiber.App, userHandler *handler.UserHandler) {
	// Test route
	router.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Test route working")
	})

	// User routes
	users := router.Group("/users")
	{
		users.Post("/", userHandler.CreateHandler)
		users.Get("/:id", userHandler.GetByIDHandler)
		users.Get("/", userHandler.GetByEmailHandler)
		users.Get("/all", userHandler.GetAllHandler)
		users.Put("/:id", userHandler.UpdateHandler)
		users.Delete("/:id", userHandler.DeleteHandler)
	}
}

// printRoutes prints all registered routes for debugging
func (app *App) printRoutes() {
	fmt.Println("Registered routes:")
	for _, route := range app.fiberApp.GetRoutes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
}

// startServer starts the Fiber server
func (app *App) startServer(port string) error {
	log.Printf("Server starting on port %s", port)
	return app.fiberApp.Listen(":" + port)
}

// waitForShutdown waits for an interrupt signal and handles graceful shutdown
func (app *App) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down server...")

	// Cancel context to stop background processes
	app.cancel()

	// Give some time for graceful shutdown
	time.Sleep(2 * time.Second)
	log.Println("Server shutdown complete")
}

// cleanup performs cleanup operations
func (app *App) cleanup() {
	if app.mongo != nil {
		app.mongo.Close()
	}
	app.cancel()
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"message": "Service is running",
		})
	}
}
