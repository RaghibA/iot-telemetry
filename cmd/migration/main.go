package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// main runs the database migrations.
// Params: None
// Returns: None
func main() {
	dbConfig, err := config.GetDBConfig()
	if err != nil {
		log.Fatal("failed to get db config", err)
	}
	dbString := config.GetDBString(dbConfig)

	db, err := sql.Open("pgx", dbString)
	if err != nil {
		log.Fatal("Failed to open db", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to init migration driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create new migrate instance", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	fmt.Println("DB Migrations complete")
}
