package main

import (
	"github.com/example/go-clean-architecture/internal/handler"
	"github.com/example/go-clean-architecture/pkg/monitoring"
	"github.com/gofiber/fiber/v2"
)

// setupRoutes wires all application routes.
func (app *App) setupRoutes() {
	app.fiberApp.Get("/health", HealthCheckHandler())
	app.fiberApp.Get("/health/memory", monitoring.MemoryHealthCheckHandler(app.memoryMonitor))

	app.fiberApp.Get("/openapi", OpenAPIDocsHandler("/openapi.json"))
	app.fiberApp.Get("/openapi.json", OpenAPISpecHandler(openAPIJSONFile, "json"))
	app.fiberApp.Get("/openapi.yaml", OpenAPISpecHandler(openAPIYAMLFile, "yaml"))

	setupUserRoutes(app.fiberApp, app.userHandler)
}

// setupUserRoutes sets up user-related routes.
func setupUserRoutes(router *fiber.App, userHandler *handler.UserHandler) {
	router.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Test route working")
	})

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

// HealthCheckHandler handles health check requests.
func HealthCheckHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"message": "Service is running",
		})
	}
}
