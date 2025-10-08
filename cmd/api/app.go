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

// App represents the application with all its components.
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

// newApp initializes and wires all application components.
func newApp() (*App, error) {
	// Create context for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize memory monitor (alerts at 80% memory usage).
	memoryMonitor := monitoring.NewMemoryMonitor(0.8)
	memoryMonitor.SetAlertHandler(func(stats monitoring.MemoryStats) {
		log.Printf("WARN: High memory usage detected - Alloc: %s, Sys: %s",
			monitoring.FormatBytes(stats.Alloc),
			monitoring.FormatBytes(stats.Sys))
	})

	// Initialize PostgreSQL database.
	db, err := driver.NewDatabase()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run auto migration for required entities.
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize MongoDB.
	mongo, err := driver.NewMongo()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Initialize repositories and use cases.
	memoryLogRepo := repository.NewMemoryLogRepository(mongo)
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize HTTP handlers.
	userHandler := handler.NewUserHandler(userUsecase)

	// Initialize Fiber app with middleware.
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())
	fiberApp.Use(monitoring.MemoryMiddleware(memoryMonitor))
	fiberApp.Use(monitoring.SimpleGoroutineMiddleware())

	// Register pprof routes for profiling.
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

// startMemoryMonitoring starts the periodic memory monitoring loop.
func (app *App) startMemoryMonitoring() {
	go app.memoryMonitor.StartMonitoring(app.ctx, 30*time.Second)
}

// startMemoryLogging starts periodic memory logging to MongoDB.
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

				// Store memory stats in MongoDB.
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

// startServer starts the Fiber HTTP server.
func (app *App) startServer(port string) error {
	log.Printf("Server starting on port %s", port)
	return app.fiberApp.Listen(":" + port)
}

// printRoutes logs all registered routes for debugging purposes.
func (app *App) printRoutes() {
	fmt.Println("Registered routes:")
	for _, route := range app.fiberApp.GetRoutes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
}

// waitForShutdown blocks until an interrupt signal is received and performs graceful shutdown.
func (app *App) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down server...")

	app.cancel()
	time.Sleep(2 * time.Second)

	log.Println("Server shutdown complete")
}

// cleanup releases all resources associated with the application.
func (app *App) cleanup() {
	if app.mongo != nil {
		app.mongo.Close()
	}
	app.cancel()
}
