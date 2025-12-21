package main

import (
	"fmt"
	"jokes-provider/api"
	"os"
)

// @title Jokes Provider API
// @version 1.0
// @description A high-performance REST API service for serving random jokes with built-in caching, rate limiting, and health checks.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Jean Yamazian
// @contact.url    https://jean.yamazian.com
// @contact.email  jeanyamazian@outlook.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Initialize the application
	app, err := api.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Graceful shutdown
	defer api.Shutdown()

	// Start the server
	if err := api.Start(app); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
