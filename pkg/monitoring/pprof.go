package monitoring

import (
	"net/http/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// RegisterPprofRoutes registers pprof routes with the Fiber application
func RegisterPprofRoutes(app *fiber.App) {
	// Create a new group for pprof endpoints
	pprofGroup := app.Group("/debug/pprof")

	// Register pprof handlers
	pprofGroup.Get("/", adaptor.HTTPHandlerFunc(pprof.Index))
	pprofGroup.Get("/cmdline", adaptor.HTTPHandlerFunc(pprof.Cmdline))
	pprofGroup.Get("/profile", adaptor.HTTPHandlerFunc(pprof.Profile))
	pprofGroup.Get("/symbol", adaptor.HTTPHandlerFunc(pprof.Symbol))
	pprofGroup.Get("/trace", adaptor.HTTPHandlerFunc(pprof.Trace))
	pprofGroup.Get("/allocs", adaptor.HTTPHandlerFunc(pprof.Handler("allocs").ServeHTTP))
	pprofGroup.Get("/block", adaptor.HTTPHandlerFunc(pprof.Handler("block").ServeHTTP))
	pprofGroup.Get("/goroutine", adaptor.HTTPHandlerFunc(pprof.Handler("goroutine").ServeHTTP))
	pprofGroup.Get("/heap", adaptor.HTTPHandlerFunc(pprof.Handler("heap").ServeHTTP))
	pprofGroup.Get("/mutex", adaptor.HTTPHandlerFunc(pprof.Handler("mutex").ServeHTTP))
	pprofGroup.Get("/threadcreate", adaptor.HTTPHandlerFunc(pprof.Handler("threadcreate").ServeHTTP))
}
