package monitoring

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MemoryStats represents memory statistics
type MemoryStats struct {
	Alloc         uint64  `json:"alloc"`         // bytes allocated and not yet freed
	TotalAlloc    uint64  `json:"totalAlloc"`    // bytes allocated (even if freed)
	Sys           uint64  `json:"sys"`           // bytes obtained from system
	NumGC         uint32  `json:"numGC"`         // number of garbage collections
	GCCPUFraction float64 `json:"gcCPUFraction"` // fraction of CPU time spent in GC
	NumGoroutine  int     `json:"numGoroutine"`  // number of goroutines
}

// MemoryMonitor represents a memory monitoring service
type MemoryMonitor struct {
	mu             sync.RWMutex
	stats          MemoryStats
	maxAlloc       uint64
	alertThreshold float64
	alertHandler   func(MemoryStats)
}

// NewMemoryMonitor creates a new memory monitor
func NewMemoryMonitor(alertThreshold float64) *MemoryMonitor {
	return &MemoryMonitor{
		alertThreshold: alertThreshold,
		maxAlloc:       0,
	}
}

// GetMemoryStats returns current memory statistics
func (m *MemoryMonitor) GetMemoryStats() MemoryStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := MemoryStats{
		Alloc:         ms.Alloc,
		TotalAlloc:    ms.TotalAlloc,
		Sys:           ms.Sys,
		NumGC:         ms.NumGC,
		GCCPUFraction: ms.GCCPUFraction,
		NumGoroutine:  runtime.NumGoroutine(),
	}

	// Update max allocation
	if ms.Alloc > m.maxAlloc {
		m.maxAlloc = ms.Alloc
	}

	// Check for memory leak alert
	if m.alertThreshold > 0 && float64(ms.Alloc) > m.alertThreshold*float64(ms.Sys) {
		if m.alertHandler != nil {
			m.alertHandler(stats)
		}
	}

	return stats
}

// GetMaxAlloc returns the maximum memory allocation observed
func (m *MemoryMonitor) GetMaxAlloc() uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxAlloc
}

// SetAlertHandler sets a callback function for memory alerts
func (m *MemoryMonitor) SetAlertHandler(handler func(MemoryStats)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alertHandler = handler
}

// StartMonitoring starts periodic memory monitoring
func (m *MemoryMonitor) StartMonitoring(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.GetMemoryStats()
		}
	}
}

// FormatBytes formats bytes into a human-readable string
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// MemoryHealthCheckHandler returns a Fiber handler for memory health checks
func MemoryHealthCheckHandler(monitor *MemoryMonitor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats := monitor.GetMemoryStats()

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "healthy",
			"memory": fiber.Map{
				"alloc":         FormatBytes(stats.Alloc),
				"totalAlloc":    FormatBytes(stats.TotalAlloc),
				"sys":           FormatBytes(stats.Sys),
				"numGC":         stats.NumGC,
				"gcCPUFraction": fmt.Sprintf("%.4f", stats.GCCPUFraction),
				"numGoroutine":  stats.NumGoroutine,
				"maxAlloc":      FormatBytes(monitor.GetMaxAlloc()),
			},
			"timestamp": time.Now().UTC(),
		})
	}
}
