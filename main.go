package main

import (
	"fmt"
	"jokes-provider/api"
	"jokes-provider/middleware"
	"os"
)

func main() {
	// Initialize the application
	app := api.Initialize()

	// Ensure Redis connection is closed on shutdown
	defer func() {
		if err := middleware.CloseRedis(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing Redis: %v\n", err)
		}
	}()

	// Start the server
	if err := api.Start(app); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
