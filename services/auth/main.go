package main

import (
    "log"

    "github.com/RaghibA/iot-telemetry/services/auth/internal/app"
)

// main is the entry point for the authentication service.
// Params: None
// Returns: None
func main() {
    log.Println("Starting auth service")
    app.Run()
}
