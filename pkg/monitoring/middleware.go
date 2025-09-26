package monitoring

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MemoryMiddleware tracks memory usage for each request
func MemoryMiddleware(monitor *MemoryMonitor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get memory stats before request
		before := monitor.GetMemoryStats()

		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get memory stats after request
		after := monitor.GetMemoryStats()

		// Calculate memory difference
		memoryDiff := int64(after.Alloc) - int64(before.Alloc)

		// Add memory usage info to response headers (optional)
		c.Response().Header.Set("X-Memory-Before", FormatBytes(before.Alloc))
		c.Response().Header.Set("X-Memory-After", FormatBytes(after.Alloc))
		c.Response().Header.Set("X-Memory-Diff", fmt.Sprintf("%+d", memoryDiff))
		c.Response().Header.Set("X-Request-Duration", duration.String())
		c.Response().Header.Set("X-Num-Goroutines", fmt.Sprintf("%d", after.NumGoroutine))

		return err
	}
}

// SimpleGoroutineMiddleware tracks goroutine count changes
func SimpleGoroutineMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record goroutine count before request
		beforeGoroutines := runtime.NumGoroutine()

		// Process request
		err := c.Next()

		// Record goroutine count after request
		afterGoroutines := runtime.NumGoroutine()

		// Add goroutine info to response headers
		c.Response().Header.Set("X-Goroutines-Before", fmt.Sprintf("%d", beforeGoroutines))
		c.Response().Header.Set("X-Goroutines-After", fmt.Sprintf("%d", afterGoroutines))
		c.Response().Header.Set("X-Goroutines-Diff", fmt.Sprintf("%d", afterGoroutines-beforeGoroutines))

		return err
	}
}
