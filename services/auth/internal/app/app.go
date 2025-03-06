package app

import (
    "log"

    "github.com/RaghibA/iot-telemetry/db"
    "github.com/RaghibA/iot-telemetry/pkg/config"
    "github.com/RaghibA/iot-telemetry/services/auth/internal/server"
)

// Run initializes the database and authentication server, and starts the server.
// Params: None
// Returns: None
func Run() {
    dbConfig, err := config.GetDBConfig()
    if err != nil {
        log.Fatal(err)
    }
    db, err := db.NewDB(dbConfig)
    if err != nil {
        log.Fatal(err)
    }

    authConfig, err := config.GetAuthConfig()
    if err != nil {
        log.Fatal(err)
    }
    s := server.NewAuthServer(authConfig, db)
    if err = s.Run(); err != nil {
        log.Fatal(err)
    }
}
