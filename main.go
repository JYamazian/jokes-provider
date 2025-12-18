package main

import (
	"fmt"
	"jokes-provider/api"
	"os"
)

func main() {
	// Initialize the application
	app := api.Initialize()

	// Start the server
	if err := api.Start(app); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
