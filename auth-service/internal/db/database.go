package db

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/auth-service/internal/models"
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
	Update(column string, conds interface{}) *gorm.DB
}

type DbInstance struct { // DbInstance holds ref to a type that implements Database interface
	Db Database
}

var IotDb DbInstance

/**
 * Migration Scripts: Users, ACLs
 *
 * Gorm Auto-migrate from models
 */

func UserMigrate() {
	err := IotDb.Db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("failed to auto-migrate user model:", err)
	}
	log.Println("User table migration complete")
}

func ACLMigrate() {
	err := IotDb.Db.AutoMigrate(&models.KafkaACL{})
	if err != nil {
		log.Println("failed to auto-migrate ACL model:", err)
	}
	log.Println("ACL table migration complete")
}

/**
 * Init DB Connection
 */
func Connect() {
	// Generate DSN from env vars
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	//! ADD RETRY LOGIC
	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Println(err)
		panic("failed to connect to postgres")
	}

	log.Println("USER DB Connected")

	// Store reference to db connection in IotDb
	IotDb = DbInstance{
		Db: db,
	}
}
