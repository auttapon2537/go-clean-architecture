package main

import "log"

func main() {
	config := loadConfig()

	app, err := newApp()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer app.cleanup()

	app.startMemoryMonitoring()
	app.startMemoryLogging()
	app.setupRoutes()
	app.printRoutes()

	if err := app.startServer(config.port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	app.waitForShutdown()
}
