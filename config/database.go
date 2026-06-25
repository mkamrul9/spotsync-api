package config

import (
	"log"
	"os"
	"time"

	"github.com/mkamrul9/spotsync-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Initialize GORM with custom logger
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Configure connection pool (Crucial for production APIs)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	sqlDB.SetMaxIdleConns(10)           // Max idle connections
	sqlDB.SetMaxOpenConns(100)          // Max open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Max lifetime of a connection

	// Run Auto-Migration
	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
	log.Println("Database migration completed successfully.")

	DB = db
	log.Println("Database connection established successfully")
}
