package db

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/iot-admin-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	AutoMigrate(...interface{}) error
	First(dest interface{}, coords ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Model(model interface{}) *gorm.DB
	Count(count *int64) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
}

type DbInstance struct { // DbInstance holds ref to a type that implements Database interface
	Db Database
}

var IotDb DbInstance

/**
 * Migrates Device model to postgres
 */
func DeviceMigrate() {
	err := IotDb.Db.AutoMigrate(&models.Device{})
	if err != nil {
		log.Println("failed to auto-migrate device model:", err)
	}
	log.Println("User table migration complete")
}

/**
 * Migrates Kafka ACL model to postgres
 */
func ACLMigrate() {
	err := IotDb.Db.AutoMigrate(&models.KafkaACL{})
	if err != nil {
		log.Println("failed to auto-migrate ACL model:", err)
	}
	log.Println("ACL table migration complete")
}

/**
 * Connects to db with env vars & stores ref to connection in IotDb
 */
func Connect() {
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	//! ADD RETRY LOGIC
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Println(err)
		panic("failed to connect to postgres")
	}

	log.Println("Admin DB Connected")

	IotDb = DbInstance{
		Db: db,
	}
}
